use std::path::Path;

use grass::OutputStyle;
use keepcalm::Shared;
use templates::Templates;

pub mod templates;

#[allow(dead_code)]
pub struct ResourceHolder {
    templates: Templates,
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

        Ok(Self {
            templates: Templates::build(&respath)?,
        })
    }
}

#[allow(dead_code)]
fn compile_css(resource_path: &Path) -> Result<String, Box<grass::Error>> {
    let opts = grass::Options::default()
        .input_syntax(grass::InputSyntax::Scss)
        .style(OutputStyle::Compressed)
        .load_path(resource_path.join("css/"));

    let out = grass::from_string("@use 'root'".to_owned(), &opts)?;

    Ok(out)
}

pub struct Resources {
    pub templates: Shared<Templates>,
}

impl Resources {}

impl From<ResourceHolder> for Resources {
    fn from(value: ResourceHolder) -> Self {
        Self {
            templates: Shared::new(value.templates),
        }
    }
}
