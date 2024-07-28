use std::{path::Path, time::Duration};

use grass::OutputStyle;
use keepcalm::{Shared, SharedMut};
use notify::Watcher;
use static_files::StaticFiles;
use templates::Templates;
use tokio::sync::watch;
use tracing::{error, info};

pub mod filters;
pub mod static_files;
pub mod templates;

struct ResourceHolder {
    templates: Templates,
    static_files: StaticFiles,
    root_static_file: StaticFiles,
}

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("Template Error: {0}")]
    Tera(#[from] tera::Error),
    #[error("Style Error")]
    Style(#[from] Box<grass::Error>),
    #[error("Directory Not Found: {0}")]
    DirNotFound(String),
    #[error("Io Error")]
    Io(#[from] std::io::Error),
    #[error("Notify Error")]
    Notify(#[from] notify::Error),
}

impl ResourceHolder {
    pub fn new<P: AsRef<Path>>(path: P) -> Result<Self, Error> {
        let path = path.as_ref();
        info!("attempting to build static directory based in: {:?}", &path);
        if !path.exists() {
            return Err(Error::DirNotFound(path.to_string_lossy().into_owned()));
        }

        let respath = path.canonicalize()?;

        let mut statics = StaticFiles::default();
        statics.register_dir(respath.join("static"))?;
        statics.register("styles", "css", compile_css(&respath)?.as_bytes());

        let mut root_statics = StaticFiles::default();
        root_statics.register_dir(respath.join("static/root"))?;

        Ok(Self {
            templates: Templates::build(&respath, statics.clone())?,
            static_files: statics,
            root_static_file: root_statics,
        })
    }
}

#[derive(Clone)]
pub struct Resources {
    pub templates: Shared<Templates>,
    pub statics: Shared<StaticFiles>,
    pub root_statics: Shared<StaticFiles>,
}

impl Resources {
    pub fn new<P: AsRef<Path>>(path: P) -> Result<Self, crate::Error> {
        Ok(ResourceHolder::new(path)?.into())
    }

    pub fn new_watched<P: AsRef<Path>>(path: P) -> Result<Self, crate::Error> {
        let path = path.as_ref();

        let r = SharedMut::new(ResourceHolder::new(path)?);

        let (tx, mut rx) = watch::channel(false);
        let mut watcher = notify::recommended_watcher(move |res| {
            if let Ok(_) = res {
                let _ = tx.send(true);
            }
        })?;

        watcher.watch(path, notify::RecursiveMode::Recursive)?;
        let path = path.to_owned();
        let r_set = r.clone();
        tokio::spawn(async move {
            while rx.changed().await.is_ok() {
                let path = path.clone();
                while tokio::time::timeout(Duration::from_millis(100), rx.changed())
                    .await
                    .is_ok()
                {}

                let res = tokio::task::spawn_blocking(move || ResourceHolder::new(path)).await;
                match res {
                    Ok(Ok(v)) => *r_set.write() = v,
                    Ok(Err(e)) => error!("Failed To Regenerate Data: {:?}", e),
                    _ => {}
                }
            }
            drop(watcher);
        });

        Ok(r.into())
    }
}

impl From<ResourceHolder> for Resources {
    fn from(value: ResourceHolder) -> Self {
        Self {
            templates: Shared::new(value.templates),
            statics: Shared::new(value.static_files),
            root_statics: Shared::new(value.root_static_file),
        }
    }
}

impl From<SharedMut<ResourceHolder>> for Resources {
    fn from(value: SharedMut<ResourceHolder>) -> Self {
        Self {
            templates: value.shared_copy().project_fn(|x| &x.templates),
            root_statics: value.shared_copy().project_fn(|x| &x.root_static_file),
            statics: value.shared_copy().project_fn(|x| &x.static_files),
        }
    }
}

fn compile_css(resource_path: &Path) -> Result<String, Box<grass::Error>> {
    let css_path = resource_path.join("css/");

    info!("compiling css in {:?}", css_path);

    let opts = grass::Options::default()
        .input_syntax(grass::InputSyntax::Scss)
        .style(OutputStyle::Compressed)
        .load_path(css_path);

    let out = grass::from_string("@use 'root'".to_owned(), &opts)?;

    Ok(out)
}
