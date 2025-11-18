import json
from typing import Any, Dict, Optional

import requests


class HTTPClient:
    def __init__(self, base_url: str, timeout: int = 30) -> None:
        self.base_url = base_url.rstrip("/") if base_url else ""
        self.timeout = timeout
        self._session = requests.Session()
        # Set default headers to avoid WAF blocking
        self._session.headers.update({
            "X-App-Token": "ANTECH-APP-SECRET-KEY-001",
            "User-Agent": "Mozilla/5.0 (Mobile; FlutterApp/1.0)",
            "Accept": "application/json",
        })

    def set_base_url(self, base_url: str) -> None:
        self.base_url = base_url.rstrip("/")

    def get(self, path: str, params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        if not self.base_url:
            raise ValueError("gateway base_url is not set")
        url = f"{self.base_url}{path}"
        full_url = url
        if params:
            from urllib.parse import urlencode
            full_url = f"{url}?{urlencode(params)}"
        
        try:
            resp = self._session.get(url, params=params or {}, timeout=self.timeout)
            
            # Check for non-JSON responses (e.g., WAF challenge pages)
            content_type = resp.headers.get('Content-Type', '').lower()
            is_html = resp.text.strip().startswith('<!DOCTYPE') or resp.text.strip().startswith('<html') or '<html' in resp.text[:100]
            
            if is_html or 'text/html' in content_type:
                # WAF challenge or HTML response
                error_info = {
                    "method": "GET",
                    "url": full_url,
                    "params": params or {},
                    "status_code": resp.status_code,
                    "content_type": content_type,
                    "response_text": resp.text[:500],
                }
                error_msg = f"Gateway returned HTML page (likely WAF challenge/CAPTCHA)\nRequest details: {error_info}"
                new_exc = requests.exceptions.HTTPError(error_msg, response=resp)
                new_exc.request_details = error_info
                raise new_exc
            
            resp.raise_for_status()
            if not resp.text:
                raise ValueError("empty response body")
            return resp.json()
        except requests.exceptions.HTTPError as e:
            # Add detailed request info for debugging
            status_code = e.response.status_code if e.response else None
            response_text = None
            if e.response:
                try:
                    response_text = e.response.text[:500]
                except Exception:
                    pass
            error_info = {
                "method": "GET",
                "url": full_url,
                "params": params or {},
                "status_code": status_code,
                "response_text": response_text,
            }
            error_msg = f"{str(e)}\nRequest details: {error_info}"
            new_exc = requests.exceptions.HTTPError(error_msg, response=e.response)
            new_exc.request_details = error_info  # Store for easy access
            raise new_exc from e
        except (json.JSONDecodeError, ValueError) as e:
            # Non-JSON response or empty response
            error_info = {
                "method": "GET",
                "url": full_url,
                "params": params or {},
                "status_code": resp.status_code if 'resp' in locals() else None,
                "response_text": resp.text[:500] if 'resp' in locals() and resp.text else None,
                "error": str(e),
            }
            raise type(e)(f"{e}\nRequest details: {error_info}") from e
        except Exception as e:
            error_info = {
                "method": "GET",
                "url": full_url,
                "params": params or {},
                "error": str(e),
            }
            raise type(e)(f"{e}\nRequest details: {error_info}") from e

    def post(self, path: str, data: Any) -> Dict[str, Any]:
        if not self.base_url:
            raise ValueError("gateway base_url is not set")
        url = f"{self.base_url}{path}"
        
        try:
            resp = self._session.post(url, json=data, timeout=self.timeout)
            resp.raise_for_status()
            if not resp.text:
                raise ValueError("empty response body")
            return resp.json()
        except requests.exceptions.HTTPError as e:
            # Add detailed request info for debugging
            status_code = e.response.status_code if e.response else None
            error_info = {
                "method": "POST",
                "url": url,
                "data": data,
                "status_code": status_code,
                "response_text": e.response.text[:500] if e.response and e.response.text else None,
            }
            error_msg = f"{str(e)}\nRequest details: {error_info}"
            new_exc = requests.exceptions.HTTPError(error_msg, response=e.response)
            new_exc.request_details = error_info  # Store for easy access
            raise new_exc from e
        except Exception as e:
            error_info = {
                "method": "POST",
                "url": url,
                "data": data,
                "error": str(e),
            }
            raise type(e)(f"{e}\nRequest details: {error_info}") from e


