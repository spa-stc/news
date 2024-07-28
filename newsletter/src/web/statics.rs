use axum::{
    body::Bytes,
    extract::{Path, State},
    http::{
        header::{CACHE_CONTROL, CONTENT_LENGTH, CONTENT_TYPE, ETAG, IF_NONE_MATCH, SERVER},
        HeaderMap, HeaderValue, StatusCode,
    },
    response::IntoResponse,
};
use lazy_static::lazy_static;
use resources::Resources;

use super::error::WebResult;

lazy_static! {
    pub static ref IMMUTABLE_CACHE_HEADER: HeaderValue = "public, max-age=31536000, immutable"
        .parse()
        .expect("Failed to parse header");
    pub static ref WELL_KNOWN_HEADER: HeaderValue = "public, max-age=86400, immutable"
        .parse()
        .expect("Failed to parse header");
    pub static ref SERVER_HEADER: HeaderValue = "spa-ssn".parse().expect("Failed to parse header");
}

#[allow(clippy::declare_interior_mutable_const)]
const NO_RESPONSE: Bytes = Bytes::new();

type StaticResponse = (StatusCode, HeaderMap, Bytes);

fn not_modified(headers: HeaderMap) -> StaticResponse {
    (StatusCode::NOT_MODIFIED, headers, NO_RESPONSE)
}

fn ok(bytes: Bytes, headers: HeaderMap) -> StaticResponse {
    (StatusCode::OK, headers, bytes)
}

fn not_found(key: &str, headers: HeaderMap) -> StaticResponse {
    tracing::warn!("File not found: {}", key);
    (StatusCode::NOT_FOUND, headers, NO_RESPONSE)
}

#[tracing::instrument(skip_all)]
pub async fn serve_statics(
    headers_in: HeaderMap,
    Path(key): Path<String>,
    State(resources): State<resources::Resources>,
) -> WebResult<impl IntoResponse> {
    let mut headers = HeaderMap::new();
    headers.append(ETAG, key.parse()?);
    headers.append(SERVER, SERVER_HEADER.clone());

    if let Some((bytes, mime)) = resources.statics.read().get_bytes_from_key(&key) {
        headers.append(CACHE_CONTROL, IMMUTABLE_CACHE_HEADER.clone());
        headers.append(CONTENT_LENGTH, bytes.len().into());
        headers.append(CONTENT_TYPE, mime.parse()?);
        if let Some(etag) = headers_in.get(IF_NONE_MATCH) {
            if *etag == key {
                return Ok(not_modified(headers));
            }
        }
        Ok(ok(bytes, headers))
    } else {
        if let Some(etag) = headers_in.get(IF_NONE_MATCH) {
            if *etag == key {
                return Ok(not_modified(headers));
            }
        }

        Ok(not_found(&key, headers))
    }
}

#[tracing::instrument(skip_all)]
pub async fn serve_root_statics(
    headers_in: HeaderMap,
    Path(file): Path<String>,
    State(resources): State<Resources>,
) -> WebResult<impl IntoResponse> {
    let mut headers = HeaderMap::new();
    headers.append(SERVER, SERVER_HEADER.clone());

    let statics = resources.root_statics.read();

    if let Some(key) = statics.lookup_key(&file) {
        headers.append(ETAG, key.parse()?);

        if let Some((bytes, mime)) = statics.get_bytes_from_key(key) {
            headers.append(CACHE_CONTROL, WELL_KNOWN_HEADER.clone());
            headers.append(CONTENT_LENGTH, bytes.len().into());
            headers.append(CONTENT_TYPE, mime.parse()?);
            if let Some(etag) = headers_in.get(IF_NONE_MATCH) {
                if *etag == key {
                    return Ok(not_modified(headers));
                }
            }
            Ok(ok(bytes, headers))
        } else {
            Ok(not_found(key, headers))
        }
    } else {
        Ok(not_found(&file, headers))
    }
}
