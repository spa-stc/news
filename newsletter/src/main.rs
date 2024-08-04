use std::{
    net::{IpAddr, SocketAddr},
    str::FromStr,
};

use eyre::bail;
use newsletter::{
    config::{self},
    web::{self, Listener},
};
use resources::Resources;
use tokio::{fs, net::TcpListener};
use tracing::info;

#[tokio::main]
async fn main() -> eyre::Result<()> {
    stable_eyre::install()?;
    let _ = kankyo::init();
    tracing_subscriber::fmt::init();

    let database = get_database().await?;

    let resources = get_res()?;

    let listener: Listener = {
        let addr = SocketAddr::new(
            IpAddr::from_str(&config::var("HOST").unwrap_or("::".into()))?,
            config::var("PORT")
                .unwrap_or("3000".into())
                .parse::<u16>()?,
        );

        info!("starting server at port: {}", addr.port());
        let tcp = TcpListener::bind(addr).await?;

        Listener::Tcp(tcp)
    };

    web::start_server(listener, resources, database).await?;

    Ok(())
}

fn get_res() -> eyre::Result<resources::Resources> {
    if let Some(dir) = config::var("PUBLIC") {
        let res = if config::is_development() {
            Resources::new_watched(&dir)?
        } else {
            Resources::new(&dir)?
        };

        Ok(res)
    } else {
        bail!("Missing NEWSLETTER_PUBLIC environment variable while fetching resources.");
    }
}

async fn get_database() -> eyre::Result<sqlx::SqlitePool> {
    if let Some(path) = config::var("DATA_DIR") {
        fs::create_dir_all(&path).await?;

        let database = store::connect(&format!("{}/data.db", &path)).await?;

        store::migrate(&database).await?;

        info!("data dir located at path: {}", path);

        return Ok(database);
    }

    bail!("Missing Required Value For DATA_DIR in environment.");
}
