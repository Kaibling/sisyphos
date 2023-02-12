
export function getURL() {
    return sessionStorage.getItem("s_url");
};

export function getToken() {
    return sessionStorage.getItem("s_token");
};

export default function InitStorage() {
    let url = import.meta.env.VITE_BASE_URL;
    sessionStorage.setItem("s_url", url)
};