use std::{borrow::Borrow, path::Path, sync::Arc};

use serde::Serialize;
use tera::{Context, Tera};
use tracing::info;

use crate::{filters::StaticFileFilter, static_files::StaticFiles};

pub struct Templates {
    tera: Tera,
}

impl Templates {
    pub fn build(resource_path: &Path, files: StaticFiles) -> Result<Self, crate::Error> {
        info!("building templates in: {:?}", resource_path);
        let mut tera = Tera::new(
            resource_path
                .join("templates/**/*")
                .to_string_lossy()
                .borrow(),
        )?;

        tera.register_filter("static", StaticFileFilter::new(files));

        Ok(Self { tera })
    }

    pub fn render<A: Serialize>(&self, name: &str, context: A) -> Result<String, crate::Error> {
        Ok(self
            .tera
            .render(name, &Context::from_serialize(&context)?)?)
    }
}

#[allow(dead_code)]
#[derive(Serialize)]
pub struct BaseRenderContext<'a, A>
where
    A: Serialize,
{
    pub title: &'a str,
    pub data: A,
    pub site_data: Arc<BaseSiteData>,
}

#[derive(Serialize)]
pub struct BaseSiteData {}
