use std::io::Read;

use chrono::{NaiveDate, Utc};
use csv::{ReaderBuilder, Terminator};
use icalendar::{
    parser::{read_calendar, unfold},
    Calendar, CalendarComponent, Component,
};
use itertools::Itertools;
use serde::Deserialize;
use store::models::days::Day;

pub struct DayScraper {
    ical: Calendar,
    csv_data: Vec<CsvData>,
}

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("Csv Error")]
    CsvError(#[from] csv::Error),

    #[error("Ical Parse Error: {0}")]
    Ical(String),
}

impl DayScraper {
    pub fn new<R: Read>(calendar: String, csv: R) -> Result<Self, Error> {
        let ical: Calendar = read_calendar(&unfold(&calendar))
            .map_err(|e| Error::Ical(e))?
            .into();

        let mut csv_reader = ReaderBuilder::new()
            .trim(csv::Trim::All)
            .terminator(Terminator::CRLF)
            .from_reader(csv);

        let csv_data: Vec<CsvData> = csv_reader
            .deserialize()
            .process_results(|iter| iter.collect())?;

        Ok(Self { ical, csv_data })
    }

    pub fn parse_days(&self, dates: Vec<NaiveDate>) -> Vec<Day> {
        struct DayData {
            date: NaiveDate,
            lunch: Option<String>,
            csvdata: Option<CsvData>,
        }

        let mut day_data: Vec<DayData> = dates
            .into_iter()
            .map(|date| DayData {
                date,
                lunch: None,
                csvdata: None,
            })
            .collect();

        for component in &self.ical.components {
            if let CalendarComponent::Event(event) = component {
                if let Some(start) = event.get_start() {
                    for day in day_data.iter_mut() {
                        if start == day.date.into() {
                            day.lunch = event.get_description().map(|e| e.into());
                        }
                    }
                }
            }
        }

        for data in &self.csv_data {
            for day in day_data.iter_mut() {
                if data.date == day.date.format("%-m/%-d/%Y").to_string() {
                    day.csvdata = Some(data.clone())
                }
            }
        }

        let days: Vec<Day> = day_data
            .into_iter()
            .map(|day| {
                let lunch = day.lunch.unwrap_or("Not Available".into());

                if let Some(d) = day.csvdata {
                    return Day {
                        date: day.date,
                        lunch: lunch,
                        ap_info: d.ap_info,
                        cc_info: d.cc_info,
                        x_period: d.x_period,
                        rotation_day: d.rotation_day,
                        location: d.location,
                        notes: "".into(),
                        grade_9: d.grade_9,
                        grade_10: d.grade_10,
                        grade_11: d.grade_11,
                        grade_12: d.grade_12,
                        updated_ts: Utc::now().naive_utc().into(),
                        created_ts: Utc::now().naive_utc().into(),
                    };
                }

                Day {
                    date: day.date,
                    lunch: lunch,
                    updated_ts: Utc::now().naive_utc().into(),
                    created_ts: Utc::now().naive_utc().into(),
                    ..Default::default()
                }
            })
            .collect();

        return days;
    }

    fn scrape_lunch(component: CalendarComponent, date: chrono::NaiveDate) -> Option<String> {
        if let CalendarComponent::Event(event) = component {
            if let Some(start) = event.get_start() {
                if start == date.into() {
                    if let Some(description) = event.get_description() {
                        return Some(description.into());
                    }
                }
            }
        }
        None
    }
}

#[derive(Deserialize, Default, Clone)]
struct CsvData {
    #[serde(alias = "DATE")]
    pub date: String,
    #[serde(alias = "EVENT")]
    pub x_period: String,
    #[serde(alias = "R. DAY")]
    pub rotation_day: String,
    #[serde(alias = "LOCATION")]
    pub location: String,
    #[serde(alias = "AP EXAMS")]
    pub ap_info: String,
    #[serde(alias = "CC TOPICS")]
    pub cc_info: String,
    #[serde(alias = "9th GRADE")]
    pub grade_9: String,
    #[serde(alias = "10th GRADE")]
    pub grade_10: String,
    #[serde(alias = "11th GRADE")]
    pub grade_11: String,
    #[serde(alias = "12th GRADE")]
    pub grade_12: String,
}

#[cfg(test)]
mod tests {
    use chrono::{NaiveDate, NaiveDateTime};
    use store::models::days::Day;

    use super::DayScraper;

    static ICAL: &'static str = include_str!("../fixtures/events.ics");
    static CSV: &'static str = include_str!("../fixtures/data.csv");

    #[test]
    fn test_parser() -> Result<(), super::Error> {
        let parser = DayScraper::new(ICAL.into(), CSV.as_bytes())?;

        let days: Vec<store::models::days::Day> = parser
            .parse_days(vec![
                NaiveDate::from_ymd_opt(2023, 8, 28).unwrap(),
                NaiveDate::from_ymd_opt(2023, 8, 29).unwrap(),
            ])
            .into_iter()
            .map(|mut d| {
                d.updated_ts = NaiveDateTime::default();
                d.created_ts = NaiveDateTime::default();
                d
            })
            .collect();

        let expected: Vec<Day> = vec![
            Day {
                date: NaiveDate::from_ymd_opt(2023, 8, 28).unwrap(),
                lunch: "Popcorn Chicken Bowl\nFresh Mashed Potato/Gravy\nRoasted Corn\nDinner roll\nPlant Based Chicken\n"
                    .to_string(),
                rotation_day: "3".into(),
                x_period: "Elected/Selected".into(),
                ..Default::default()
            },
            Day {
                date: NaiveDate::from_ymd_opt(2023, 8, 29).unwrap(),
                x_period: "Open, Sr. Meeting".into(),
                grade_12: "Sr. Speeches?".into(),
                lunch: "Cheese Enchilada's\nChicken Tinga\nLime Cilantro Rice\nCauliflower and Roasted Tomatoes\nCinnamon Churro's\n\n"
                    .into(),
                rotation_day: "4".into(),
                ..Default::default()
            },
        ];
        assert_eq!(days, expected);

        Ok(())
    }
}
