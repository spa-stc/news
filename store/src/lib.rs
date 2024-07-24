use std::str::FromStr;

use sqlx::{error::ErrorKind, migrate, sqlite::SqliteConnectOptions, SqlitePool};

pub mod models;

pub type Result<T, E = Error> = std::result::Result<T, E>;

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("Record Not Found")]
    NotFound,
    #[error("Insertion Conflict")]
    Conflict,
    #[error("Sqlx Error: {0}")]
    Sqlx(sqlx::Error),
}

impl From<sqlx::Error> for Error {
    fn from(err: sqlx::Error) -> Self {
        match err {
            sqlx::Error::RowNotFound => crate::Error::NotFound,
            sqlx::Error::Database(dbe) if dbe.kind() == ErrorKind::UniqueViolation => {
                crate::Error::Conflict
            }
            e => Error::Sqlx(e),
        }
    }
}

pub async fn connect(database_path: &str) -> Result<SqlitePool, sqlx::Error> {
    let conn_ops = SqliteConnectOptions::from_str(database_path)?
        .journal_mode(sqlx::sqlite::SqliteJournalMode::Wal)
        .create_if_missing(true)
        .pragma("foreign_keys", "0")
        .pragma("busy_timeout", "1000");

    SqlitePool::connect_with(conn_ops).await
}

pub async fn connect_inmem() -> Result<SqlitePool, sqlx::Error> {
    let conn_ops = SqliteConnectOptions::from_str("sqlite::memory:")?
        .pragma("foreign_keys", "0")
        .pragma("busy_timeout", "1000");

    SqlitePool::connect_with(conn_ops).await
}

pub async fn migrate(conn: &SqlitePool) -> Result<(), sqlx::Error> {
    migrate!().run(conn).await?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_migrations() -> sqlx::Result<()> {
        let pool = connect_inmem().await?;

        migrate(&pool).await?;

        Ok(())
    }
}
