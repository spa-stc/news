use sqlx::error::ErrorKind;

pub fn get_dberr(err: sqlx::Error) -> crate::Error {
    match err {
        sqlx::Error::RowNotFound => crate::Error::NotFound,
        sqlx::Error::Database(dbe) if dbe.kind() == ErrorKind::UniqueViolation => {
            crate::Error::Conflict
        }
        _ => err.into(),
    }
}
