use sqlx::SqliteExecutor;

use crate::utils::get_dberr;

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord)]
pub struct NewUser {
    pub email: String,
    pub name: String,
    pub password: String,
    pub is_admin: bool,
}

#[derive(Debug, Clone, Default, PartialEq, Eq, PartialOrd, Ord)]
pub struct UpdateUser {
    pub name: Option<String>,
    pub password: Option<String>,
    pub is_admin: Option<bool>,
    pub email_verified: Option<bool>,
    pub disabled: Option<bool>,
}

#[derive(Debug, sqlx::FromRow, Clone, Default, PartialEq, Eq, PartialOrd, Ord)]
pub struct User {
    pub id: i32,
    pub email: String,
    pub name: String,
    pub password: String,
    pub is_admin: bool,
    pub disabled: bool,
    pub email_verified: bool,
    pub created_ts: chrono::NaiveDateTime,
    pub updated_ts: chrono::NaiveDateTime,
}

impl User {
    pub async fn insert<'a>(
        executor: impl SqliteExecutor<'a>,
        new: NewUser,
    ) -> crate::Result<Self> {
        let row: Self = sqlx::query_as(
            r#"
        INSERT INTO users (
           email, 
           name, 
           password, 
           is_admin
        ) 
        VALUES (
            ?1, 
            ?2, 
            ?3, 
            ?4
        )
        RETURNING
            id, 
            email, 
            name, 
            password,
            is_admin, 
            disabled,
            email_verified,
            created_ts, 
            updated_ts
            "#,
        )
        .bind(&new.email)
        .bind(&new.name)
        .bind(&new.password)
        .bind(&new.is_admin)
        .fetch_one(executor)
        .await
        .map_err(get_dberr)?;

        Ok(row)
    }

    pub async fn get<'a>(executor: impl SqliteExecutor<'a>, id: i32) -> crate::Result<Self> {
        let row: Self = sqlx::query_as(
            r#"
        SELECT 
            id, 
            email, 
            name, 
            password,
            is_admin,
            email_verified,
            disabled,
            created_ts, 
            updated_ts
        FROM 
            users
        WHERE
            id = ?1
        "#,
        )
        .bind(id)
        .fetch_one(executor)
        .await
        .map_err(get_dberr)?;

        Ok(row)
    }

    pub async fn update<'a>(
        executor: impl SqliteExecutor<'a>,
        usr: User,
        update: UpdateUser,
    ) -> crate::Result<Self> {
        let row: Self = sqlx::query_as(
            r#"
        UPDATE 
            users
        SET 
            name = ?2, 
            password = ?3, 
            is_admin = ?4,
            email_verified = ?5,
            disabled = ?6
        WHERE
            id = ?1
        RETURNING
            id, 
            email, 
            name, 
            password,
            is_admin,
            email_verified,
            disabled,
            created_ts, 
            updated_ts
        "#,
        )
        .bind(&usr.id)
        .bind(&update.name.unwrap_or(usr.name))
        .bind(&update.password.unwrap_or(usr.password))
        .bind(&update.is_admin.unwrap_or(usr.is_admin))
        .bind(&update.email_verified.unwrap_or(usr.email_verified))
        .bind(&update.disabled.unwrap_or(usr.disabled))
        .fetch_one(executor)
        .await
        .map_err(get_dberr)?;

        Ok(row)
    }
}

#[cfg(test)]
mod tests {
    use chrono::DateTime;
    use sqlx::SqlitePool;

    use crate::{connect_inmem, migrate, Error};

    use super::*;

    async fn get_store() -> crate::Result<SqlitePool> {
        let store = connect_inmem().await?;
        migrate(&store).await?;

        Ok(store)
    }

    #[tokio::test]
    async fn test_user_insert_get() -> crate::Result<()> {
        let store = get_store().await?;

        let nu = NewUser {
            email: "e@e.com".into(),
            name: "Foehammer".into(),
            password: "1234567".into(),
            is_admin: false,
        };

        let new_user = User::insert(&store, nu.clone()).await?;

        assert_eq!(new_user.email, nu.email);
        assert_eq!(new_user.name, nu.name);
        assert_eq!(new_user.password, nu.password);
        assert_eq!(new_user.is_admin, nu.is_admin);

        let user = User::get(&store, new_user.id).await?;

        assert_eq!(user, new_user);

        Ok(())
    }

    #[tokio::test]
    async fn test_user_insert_failure() -> crate::Result<()> {
        let store = get_store().await?;

        let nu = NewUser {
            email: "e@e.com".into(),
            name: "Foehammer".into(),
            password: "1234567".into(),
            is_admin: false,
        };

        let _ = User::insert(&store, nu.clone()).await?;

        match User::insert(&store, nu).await.unwrap_err() {
            Error::Conflict => {}
            err => {
                println!("{:?}", err);
                assert!(false);
            }
        }

        Ok(())
    }

    #[tokio::test]
    async fn test_user_get_failiure() -> crate::Result<()> {
        let pool = get_store().await?;

        match User::get(&pool, 1234).await.unwrap_err() {
            Error::NotFound => {}
            e => {
                println!("{:?}", e);
                assert!(false);
            }
        }

        Ok(())
    }

    #[tokio::test]
    async fn test_user_update() -> crate::Result<()> {
        let store = get_store().await?;

        let nu = NewUser {
            email: "e@e.com".into(),
            name: "Foehammer".into(),
            password: "1234567".into(),
            is_admin: false,
        };

        let usr = User::insert(&store, nu.clone()).await?;

        let update_user = UpdateUser {
            name: Some("hello".into()),
            password: Some("hi".into()),
            is_admin: Some(true),
            email_verified: Some(true),
            disabled: Some(true),
            ..Default::default()
        };

        let user = User::update(&store, usr.clone(), update_user.clone()).await?;

        assert_ne!(user, usr);

        Ok(())
    }

    #[tokio::test]
    async fn test_user_update_failure() -> crate::Result<()> {
        let store = get_store().await?;

        let update_user = UpdateUser {
            name: Some("hello".into()),
            password: Some("hi".into()),
            is_admin: Some(true),
            email_verified: Some(true),
            disabled: Some(true),
            ..Default::default()
        };

        let usr = User {
            id: 12,
            email: "hi".into(),
            name: "hi".into(),
            password: "hi".into(),
            is_admin: false,
            disabled: true,
            email_verified: false,
            created_ts: DateTime::from_timestamp(120, 12).unwrap().naive_utc(),
            updated_ts: DateTime::from_timestamp(120, 12).unwrap().naive_utc(),
        };

        match User::update(&store, usr, update_user).await.unwrap_err() {
            Error::NotFound => {}
            e => {
                println!("{:?}", e);
                assert!(false);
            }
        };

        Ok(())
    }
}
