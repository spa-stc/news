use axum::{extract::State, response::IntoResponse, routing::get, Router};
use resources::Resources;
use sqlx::SqlitePool;
use store::models::days::Day;
use tracing::instrument;

use super::{
    error::WebResult,
    templates::{PageCachePolicy, TemplateRenderer},
};

#[derive(Clone)]
struct PublicState {
    pub templates: TemplateRenderer,
    pub database: SqlitePool,
}

pub fn routes<S>(resources: Resources, database: SqlitePool) -> Router<S>
where
    S: Clone + Send + Sync + 'static,
{
    let state = PublicState {
        templates: TemplateRenderer::new(resources),
        database,
    };

    Router::new().route("/", get(index)).with_state(state)
}

#[derive(serde::Serialize)]
struct IndexData {
    pub days: Vec<Day>,
}

#[instrument(skip_all)]
async fn index(State(state): State<PublicState>) -> WebResult<impl IntoResponse> {
    let days = Day::get_many(&state.database, vec![]).await?;

    state.templates.render(
        PageCachePolicy::Public,
        "index.html",
        "Home",
        IndexData { days },
    )
}
