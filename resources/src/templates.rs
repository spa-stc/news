use std::{borrow::Borrow, path::Path};

use serde::Serialize;
use tera::{Context, Tera};

pub struct Templates {
    tera: Tera,
}

impl Templates {
    pub fn build(resource_path: &Path) -> Result<Self, crate::Error> {
        let tera = Tera::new(
            resource_path
                .join("templates/**/*")
                .to_string_lossy()
                .borrow(),
        )?;

        Ok(Self { tera })
    }

    pub fn render<A: Serialize>(
        &self,
        name: &str,
        context: BaseRenderContext<A>,
    ) -> Result<String, crate::Error> {
        Ok(self
            .tera
            .render(name, &Context::from_serialize(&context)?)?)
    }

    pub fn render_base<A: Serialize>(
        &self,
        name: &str,
        context: A,
    ) -> Result<String, crate::Error> {
        Ok(self
            .tera
            .render(name, &Context::from_serialize(&context)?)?)
    }
}

#[derive(Serialize)]
pub struct BaseRenderContext<'a, A>
where
    A: Serialize,
{
    title: &'a str,
    data: A,
}
