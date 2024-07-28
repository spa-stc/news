use axum::{extract::State, response::IntoResponse, routing::get, Router};
use resources::Resources;
use tracing::instrument;

use super::{
    error::WebResult,
    templates::{PageCachePolicy, TemplateRenderer},
};

pub fn routes<S>(resources: Resources) -> Router<S>
where
    S: Clone + Send + Sync + 'static,
{
    Router::new()
        .route("/", get(index))
        .with_state(TemplateRenderer::new(resources))
}

#[instrument(skip_all)]
async fn index(State(templates): State<TemplateRenderer>) -> WebResult<impl IntoResponse> {
    templates.render(PageCachePolicy::Public, "index.html", "Home", {})
}
