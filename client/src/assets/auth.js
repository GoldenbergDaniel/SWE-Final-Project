import { writable } from 'svelte/store';

export const isAuthenticated = writable(false);

export function checkAuth() {
    const sessionCookie = document.cookie.split(';').find(c => c.trim().startsWith('session_token='));
    isAuthenticated.set(!!sessionCookie);
    return isAuthenticated
}

// export async function login(username, password) {
//     const response = await fetch('http://localhost:5174/login', {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/json',
//         },
//         body: JSON.stringify({ username, password }),
//         credentials: 'include',
//     });

//     if (response.ok) {
//         isAuthenticated.set(true);
//         return true;
//     } else {
//         throw new Error('Login failed');
//     }
// }

// export async function logout() {
//     const response = await fetch('http://localhost:5174/logout', {
//         method: 'POST',
//         credentials: 'include',
//     });

//     if (response.ok) {
//         isAuthenticated.set(false);
//     } else {
//         throw new Error('Logout failed');
//     }
// }
