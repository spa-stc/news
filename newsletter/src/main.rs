use std::{
    net::{IpAddr, SocketAddr},
    str::FromStr,
};

use newsletter::{
    config,
    web::{self, Listener},
};
use tokio::net::TcpListener;

#[tokio::main]
async fn main() -> color_eyre::Result<()> {
    color_eyre::install()?;
    let _ = kankyo::init();
    tracing_subscriber::fmt::init();

    let listener: Listener = {
        let addr = SocketAddr::new(
            IpAddr::from_str(&config::var("HOST").unwrap_or("::".into()))?,
            config::var("PORT")
                .unwrap_or("3000".into())
                .parse::<u16>()?,
        );

        let tcp = TcpListener::bind(addr).await?;

        Listener::Tcp(tcp)
    };

    web::start_server(listener).await?;

    Ok(())
}
