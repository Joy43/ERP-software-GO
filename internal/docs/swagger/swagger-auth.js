(function() {
    const TOKEN_KEY = 'swagger_access_token';
    const REFRESH_TOKEN_KEY = 'swagger_refresh_token';

    function parseJwt(token) {
        try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
                return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
            }).join(''));
            return JSON.parse(jsonPayload);
        } catch (e) {
            return null;
        }
    }

    function isExpired(token) {
        const decoded = parseJwt(token);
        if (!decoded || !decoded.exp) return true;
        // Check if expiring in the next 10 seconds
        return (decoded.exp * 1000) < (Date.now() + 10000);
    }

    async function refreshToken() {
        const refresh = localStorage.getItem(REFRESH_TOKEN_KEY);
        if (!refresh) return null;

        try {
            const response = await fetch('/api/v1/auth/refresh', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refresh })
            });
            const result = await response.json();
            if (result.success && result.data && result.data.access_token) {
                localStorage.setItem(TOKEN_KEY, result.data.access_token);
                if (result.data.refresh_token) {
                    localStorage.setItem(REFRESH_TOKEN_KEY, result.data.refresh_token);
                }
                return result.data.access_token;
            }
        } catch (e) {
            console.error('Token refresh failed', e);
        }
        
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(REFRESH_TOKEN_KEY);
        return null;
    }

    // Capture SwaggerUIBundle call to inject interceptors
    const originalSwaggerUIBundle = window.SwaggerUIBundle;
    window.SwaggerUIBundle = function(config) {
        const originalRequestInterceptor = config.requestInterceptor || (req => req);
        const originalResponseInterceptor = config.responseInterceptor || (res => res);

        config.requestInterceptor = async (req) => {
            let token = localStorage.getItem(TOKEN_KEY);
            
            if (!req.url.endsWith('/auth/login') && !req.url.endsWith('/auth/refresh')) {
                if (token) {
                    if (isExpired(token)) {
                        token = await refreshToken();
                    }
                    if (token) {
                        req.headers['Authorization'] = 'Bearer ' + token;
                    }
                }
            }
            return originalRequestInterceptor(req);
        };

        config.responseInterceptor = (res) => {
            if (res.url.endsWith('/auth/login') && res.status === 200) {
                const data = res.body && res.body.data;
                if (data && data.access_token) {
                    localStorage.setItem(TOKEN_KEY, data.access_token);
                    if (data.refresh_token) {
                        localStorage.setItem(REFRESH_TOKEN_KEY, data.refresh_token);
                    }
                }
            }
            return originalResponseInterceptor(res);
        };

        return originalSwaggerUIBundle(config);
    };

    // Copy static methods from original bundle
    Object.assign(window.SwaggerUIBundle, originalSwaggerUIBundle);

})();
