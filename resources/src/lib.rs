use std::path::Path;

use grass::OutputStyle;
use keepcalm::Shared;
use static_files::StaticFiles;
use templates::Templates;

pub mod static_files;
pub mod templates;

#[allow(dead_code)]
pub struct ResourceHolder {
    templates: Templates,
    static_files: StaticFiles,
    root_static_file: StaticFiles,
}

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("Template Error")]
    Tera(#[from] tera::Error),
    #[error("Style Error")]
    Style(#[from] Box<grass::Error>),
    #[error("Directory Not Found: {0}")]
    DirNotFound(String),
    #[error("Io Error")]
    Io(#[from] std::io::Error),
}

impl ResourceHolder {
    pub fn new<P: AsRef<Path>>(path: P) -> Result<Self, Error> {
        let path = path.as_ref();
        if !path.exists() {
            return Err(Error::DirNotFound(path.to_string_lossy().into_owned()));
        }

        let respath = path.canonicalize()?;

        let mut statics = StaticFiles::default();
        statics.register_dir(respath.join("/static"))?;
        statics.register("styles", "css", compile_css(&respath)?.as_bytes());

        let mut root_statics = StaticFiles::default();
        root_statics.register_dir(respath.join("/static/root"))?;

        Ok(Self {
            templates: Templates::build(&respath)?,
            static_files: statics,
            root_static_file: root_statics,
        })
    }
}

pub struct Resources {
    pub templates: Shared<Templates>,
    pub statics: Shared<StaticFiles>,
    pub root_statics: Shared<StaticFiles>,
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

fn compile_css(resource_path: &Path) -> Result<String, Box<grass::Error>> {
    let opts = grass::Options::default()
        .input_syntax(grass::InputSyntax::Scss)
        .style(OutputStyle::Compressed)
        .load_path(resource_path.join("css/"));

    let out = grass::from_string("@use 'root'".to_owned(), &opts)?;

    Ok(out)
}
