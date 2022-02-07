
export function set(key, value) {
    localStorage.setItem(key, JSON.stringify(value))
}

export function get(key) {
    try {
        const v = localStorage.getItem(key)
        return v ? JSON.parse(v) : undefined
    } catch {
        return undefined
    }
}

export function remove(key) {
    localStorage.removeItem(key)
}
