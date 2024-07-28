use axum::{routing::get, Router};
use tokio::net::TcpListener;

pub enum Listener {
    Tcp(TcpListener),
}

#[allow(dead_code)]
pub async fn start_server(listener: Listener) -> Result<(), eyre::Report> {
    let app = Router::new().route("/healthz", get(healthz));

    match listener {
        Listener::Tcp(l) => axum::serve(l, app.into_make_service()).await?,
    };

    Ok(())
}

async fn healthz() -> &'static str {
    "Service Ready."
}
