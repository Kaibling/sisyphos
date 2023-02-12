export function ValOrNull(data) {
    if (data == null) {
        return ""
    }
    return data
}

export function GetValue(obj) {
    if (Array.isArray(obj)) {
        return JSON.stringify(obj); 
    } else if (typeof obj == "object") {
        return JSON.stringify(obj); 
    }
    return obj
}


export function ToObj(obj) {
    if (typeof obj == "object") {
        return obj
    }
    return JSON.parse(obj)
}


