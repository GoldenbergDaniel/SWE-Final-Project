import { writable } from 'svelte/store';

export const isAuthenticated = writable(false);

export function checkAuth()
{
    const sessionCookie = document.cookie.split(';').find(c => c.trim().startsWith('session_token='));
    isAuthenticated.set(sessionCookie != undefined)
    return sessionCookie != undefined
}
