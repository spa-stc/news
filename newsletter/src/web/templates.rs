use axum::{
    http::{header::CACHE_CONTROL, StatusCode},
    response::{AppendHeaders, Html, IntoResponse},
};
use resources::{templates::BaseRenderContext, Resources};
use serde::Serialize;

use super::error::WebResult;

#[derive(Clone)]
pub struct TemplateRenderer {
    resources: resources::Resources,
}

#[allow(dead_code)]
pub enum PageCachePolicy {
    Public,
    Private,
    Other(&'static str),
}

impl PageCachePolicy {
    pub fn get(&self) -> &'static str {
        match self {
            Self::Other(s) => *s,
            PageCachePolicy::Public => "public, max-age=3600, stale-if-error=60",
            PageCachePolicy::Private => "no-store",
        }
    }
}

impl TemplateRenderer {
    pub fn new(r: Resources) -> Self {
        Self { resources: r }
    }

    pub fn render<'a, A: Serialize + 'a>(
        &self,
        cache_policy: PageCachePolicy,
        name: &'a str,
        title: &'a str,
        data: A,
    ) -> WebResult<impl IntoResponse + 'a> {
        let base_context = BaseRenderContext { title, data };

        self.render_base(cache_policy, name, base_context)
    }

    pub fn render_base<A: Serialize>(
        &self,
        cache_policy: PageCachePolicy,
        name: &str,
        data: A,
    ) -> WebResult<impl IntoResponse> {
        let tmpl = self.resources.templates.read();

        Ok((
            StatusCode::OK,
            AppendHeaders([(CACHE_CONTROL, cache_policy.get())]),
            Html(tmpl.render(name, data)?),
        ))
    }
}
