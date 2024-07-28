use axum::{
    http::{header::InvalidHeaderValue, StatusCode},
    response::IntoResponse,
};
use tracing::error;

pub type WebResult<T> = Result<T, WebError>;

#[derive(Debug, thiserror::Error)]
pub enum WebError {
    #[error("Store Error")]
    Store(#[from] store::Error),

    #[error("Invalid Header Value")]
    InvalidHeaderValue(#[from] InvalidHeaderValue),

    #[error("Resource Error: {0}")]
    Resource(#[from] resources::Error),

    #[error(transparent)]
    Other(#[from] eyre::Report),
}

impl IntoResponse for WebError {
    fn into_response(self) -> axum::response::Response {
        error!("Web Error: {}", self);

        match self {
            Self::Store(store::Error::NotFound) => {
                (StatusCode::NOT_FOUND, "Not Found.").into_response()
            }
            _ => (StatusCode::INTERNAL_SERVER_ERROR, "Internal Server Error.").into_response(),
        }
    }
}
