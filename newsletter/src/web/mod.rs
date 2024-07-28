use axum::{routing::get, Router};
use resources::Resources;
use statics::{serve_root_statics, serve_statics};
use tokio::net::TcpListener;

mod error;
mod statics;

pub enum Listener {
    Tcp(TcpListener),
}

#[allow(dead_code)]
pub async fn start_server(listener: Listener, resources: Resources) -> Result<(), eyre::Report> {
    let app = Router::new()
        .route("/healthz", get(healthz))
        .route("/static/:file", get(serve_statics))
        .route("/:file", get(serve_root_statics))
        .with_state(resources);

    match listener {
        Listener::Tcp(l) => axum::serve(l, app.into_make_service()).await?,
    };

    Ok(())
}

async fn healthz() -> &'static str {
    "Service Ready."
}
