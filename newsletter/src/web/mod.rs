use axum::{
    extract::Request,
    http::{header, StatusCode},
    middleware::{self, Next},
    response::Response,
    routing::get,
};
use resources::Resources;
use statics::{serve_root_statics, serve_statics};
use tokio::net::TcpListener;

mod error;
mod public;
mod statics;
mod templates;

pub enum Listener {
    Tcp(TcpListener),
}

#[allow(dead_code)]
pub async fn start_server(listener: Listener, resources: Resources) -> Result<(), eyre::Report> {
    let app = public::routes(resources.clone())
        .route("/healthz", get(healthz))
        .route("/static/:file", get(serve_statics))
        .route("/:file", get(serve_root_statics))
        .route_layer(middleware::from_fn(request_trace))
        .with_state(resources.clone());

    match listener {
        Listener::Tcp(l) => {
            axum::serve(l, app.into_make_service()).await?;
        }
    };

    Ok(())
}

async fn healthz() -> &'static str {
    "Service Ready."
}

async fn request_trace(req: Request, next: Next) -> Result<Response, StatusCode> {
    let uri = req.uri().to_string();

    let ua = req
        .headers()
        .get(header::USER_AGENT)
        .map(|s| String::from_utf8_lossy(s.as_bytes()));

    let ip = req
        .headers()
        .get("x-forwarded-for")
        .map(|s| String::from_utf8_lossy(s.as_bytes()));

    let r = req
        .headers()
        .get(header::REFERER)
        .map(|s| String::from_utf8_lossy(s.as_bytes()));

    tracing::info!(
        "http_request {}",
        serde_json::json!({ "uri": uri, "ua": ua, "ip": ip, "r": r })
    );

    Ok(next.run(req).await)
}
