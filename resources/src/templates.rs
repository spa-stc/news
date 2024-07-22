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
        Ok(self.tera.render(name, &context.into())?)
    }

    pub fn render_base<A: Serialize>(
        &self,
        name: &str,
        context: A,
    ) -> Result<String, crate::Error> {
        let mut render_context = Context::new();

        render_context.insert("data", &context);

        Ok(self.tera.render(name, &render_context)?)
    }
}

pub struct BaseRenderContext<'a, A>
where
    A: Serialize,
{
    title: &'a str,
    data: A,
}

impl<'a, A: Serialize> Into<tera::Context> for BaseRenderContext<'a, A> {
    fn into(self) -> tera::Context {
        let mut context = Context::new();

        context.insert("title", self.title);
        context.insert("data", &self.data);

        context
    }
}
