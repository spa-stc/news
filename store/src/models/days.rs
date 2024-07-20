use chrono::NaiveDate;
use serde::{Deserialize, Serialize};
use sqlx::{sqlite::SqliteArguments, Arguments, SqliteExecutor};

use crate::utils::get_dberr;

// Day model, representing a day serialized to the database.
#[derive(Debug, sqlx::FromRow, Default, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub struct Day {
    pub date: chrono::NaiveDate,
    pub lunch: String,
    pub x_period: String,
    pub rotation_day: String,
    pub location: String,
    pub notes: String,
    pub ap_info: String,
    pub cc_info: String,
    pub grade_9: String,
    pub grade_10: String,
    pub grade_11: String,
    pub grade_12: String,
    pub created_ts: chrono::NaiveDateTime,
    pub updated_ts: chrono::NaiveDateTime,
}

impl Day {
    pub async fn upsert<'a>(&self, executor: impl SqliteExecutor<'a>) -> crate::Result<()> {
        sqlx::query(
            r#"
        INSERT INTO days (
            date,
            lunch, 
            x_period, 
            rotation_day, 
            location, 
            notes,
            ap_info,
            cc_info, 
            grade_9, 
            grade_10, 
            grade_11,
            grade_12,
            created_ts,
            updated_ts
        ) 
        VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?14) 
        ON CONFLICT DO UPDATE 
        SET
            date = ?1, 
            lunch = ?2,
            x_period = ?3,
            rotation_day = ?4,
            location = ?5,
            notes = ?6,
            ap_info = ?7,
            cc_info = ?8,
            grade_9 = ?9,
            grade_10 = ?10,
            grade_11 = ?11,
            grade_12 = ?12,
            updated_ts = ?14 
		"#,
        )
        .bind(&self.date)
        .bind(&self.lunch)
        .bind(&self.x_period)
        .bind(&self.rotation_day)
        .bind(&self.location)
        .bind(&self.notes)
        .bind(&self.ap_info)
        .bind(&self.cc_info)
        .bind(&self.grade_9)
        .bind(&self.grade_10)
        .bind(&self.grade_11)
        .bind(&self.grade_12)
        .bind(&self.created_ts)
        .bind(&self.updated_ts)
        .execute(executor)
        .await
        .map_err(get_dberr)?;

        Ok(())
    }

    pub async fn get<'a>(
        executor: impl SqliteExecutor<'a>,
        day: &NaiveDate,
    ) -> crate::Result<Self> {
        let day: Self = sqlx::query_as(
            r#"
        SELECT
            date,
            lunch, 
            x_period, 
            rotation_day, 
            location, 
            notes,
            ap_info,
            cc_info, 
            grade_9, 
            grade_10, 
            grade_11,
            grade_12,
            created_ts,
            updated_ts
        FROM 
            days 
        WHERE
            date = ?1 
            "#,
        )
        .bind(&day)
        .fetch_one(executor)
        .await
        .map_err(get_dberr)?;

        Ok(day)
    }

    pub async fn get_many<'a>(
        executor: impl SqliteExecutor<'a>,
        dates: Vec<NaiveDate>,
    ) -> crate::Result<Vec<Day>> {
        let sql = format!(
            r#"
        SELECT
            date,
            lunch, 
            x_period, 
            rotation_day, 
            location, 
            notes,
            ap_info,
            cc_info, 
            grade_9, 
            grade_10, 
            grade_11,
            grade_12,
            created_ts,
            updated_ts
        FROM 
            days 
        WHERE
            date IN ({})
        "#,
            (0..dates.len())
                .into_iter()
                .map(|_| "?")
                .collect::<Vec<&str>>()
                .join(",")
        );

        let mut args = SqliteArguments::default();
        dates
            .into_iter()
            .map(|date| {
                args.add(date);
            })
            .count();

        let days: Vec<Day> = sqlx::query_as_with(&sql, args)
            .fetch_all(executor)
            .await
            .map_err(get_dberr)?;

        if days.len() == 0 {
            return Err(crate::Error::NotFound);
        }

        Ok(days)
    }
}

#[cfg(test)]
mod test {
    use chrono::{DateTime, NaiveDate};
    use sqlx::SqlitePool;

    use crate::{connect_inmem, migrate, Error};

    use super::*;

    async fn get_store() -> crate::Result<SqlitePool> {
        let store = connect_inmem().await?;
        migrate(&store).await?;

        Ok(store)
    }

    #[tokio::test]
    async fn test_days_upsert() -> crate::Result<()> {
        let pool = get_store().await?;

        let mut day = Day {
            date: NaiveDate::parse_from_str("2006-12-12", "%Y-%m-%d").unwrap(),
            lunch: String::from("Hamburgers and Cheese"),
            x_period: String::from("Sus"),
            rotation_day: String::from(""),
            location: "e".into(),
            notes: "e".into(),
            ap_info: "e".into(),
            cc_info: "e".into(),
            grade_9: "e".into(),
            grade_10: "e".into(),
            grade_11: "e".into(),
            grade_12: "e".into(),
            created_ts: DateTime::parse_from_rfc3339("2024-07-16 21:01:28-05:00")
                .unwrap()
                .naive_utc(),
            updated_ts: DateTime::parse_from_rfc3339("2024-07-16 21:01:28-05:00")
                .unwrap()
                .naive_utc(),
        };

        day.upsert(&pool).await?;

        let new = Day::get(&pool, &day.date).await?;

        assert_eq!(day, new);

        // Ensure that the created time is not updated.
        day.created_ts = DateTime::from_timestamp(120, 12).unwrap().naive_utc();

        day.upsert(&pool).await?;

        let new2 = Day::get(&pool, &day.date).await?;

        assert_ne!(day, new2);

        Ok(())
    }

    #[tokio::test]
    async fn test_days_query_failure() -> crate::Result<()> {
        let pool = get_store().await?;

        let mut tx = pool.begin().await?;

        let _ = match Day::get(&mut *tx, &NaiveDate::from_ymd_opt(2023, 12, 2).unwrap()).await {
            Ok(_) => assert!(false),
            Err(crate::Error::NotFound) => {}
            Err(_) => {
                assert!(false);
            }
        };

        tx.commit().await?;
        Ok(())
    }

    #[tokio::test]
    async fn test_days_get_many() -> crate::Result<()> {
        let pool = get_store().await?;

        let mut tx = pool.begin().await?;

        let day = Day {
            date: NaiveDate::parse_from_str("2006-12-12", "%Y-%m-%d").unwrap(),
            lunch: String::from("Hamburgers and Cheese"),
            x_period: String::from("Sus"),
            rotation_day: String::from(""),
            location: "e".into(),
            notes: "e".into(),
            ap_info: "e".into(),
            cc_info: "e".into(),
            grade_9: "e".into(),
            grade_10: "e".into(),
            grade_11: "e".into(),
            grade_12: "e".into(),
            created_ts: DateTime::parse_from_rfc3339("2024-07-16 21:01:28-05:00")
                .unwrap()
                .naive_utc(),
            updated_ts: DateTime::parse_from_rfc3339("2024-07-16 21:01:28-05:00")
                .unwrap()
                .naive_utc(),
        };

        day.upsert(&mut *tx).await?;

        let day2 = Day {
            date: NaiveDate::parse_from_str("2006-12-13", "%Y-%m-%d").unwrap(),
            lunch: String::from("Hamburgers and Cheese"),
            x_period: String::from("Sus"),
            rotation_day: String::from(""),
            location: "e".into(),
            notes: "e".into(),
            ap_info: "e".into(),
            cc_info: "e".into(),
            grade_9: "e".into(),
            grade_10: "e".into(),
            grade_11: "e".into(),
            grade_12: "e".into(),
            created_ts: DateTime::parse_from_rfc3339("2024-07-16 21:01:28-05:00")
                .unwrap()
                .naive_utc(),
            updated_ts: DateTime::parse_from_rfc3339("2024-07-16 21:01:28-05:00")
                .unwrap()
                .naive_utc(),
        };

        day2.upsert(&mut *tx).await?;

        let days = Day::get_many(&mut *tx, vec![day.date, day2.date]).await?;
        let expected = vec![day, day2];

        tx.commit().await?;

        assert_eq!(days, expected);

        Ok(())
    }

    #[tokio::test]
    async fn test_days_get_many_failure() -> crate::Result<()> {
        let pool = get_store().await?;

        match Day::get_many(
            &pool,
            vec![NaiveDate::parse_from_str("2006-12-13", "%Y-%m-%d").unwrap()],
        )
        .await
        .unwrap_err()
        {
            Error::NotFound => {}
            e => {
                println!("{:?}", e);
                assert!(false);
            }
        }

        Ok(())
    }
}
