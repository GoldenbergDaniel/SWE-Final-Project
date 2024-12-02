<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { checkAuth } from '../auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    let posts = [];

    onMount(async () => {
        if (!checkAuth()) {
            alert('Access Denied! Login Required');
            navigate('/');
        } else {
            await fetchPosts();
        }
    });

    async function fetchPosts() {
        try {
            const response = await fetch("http://localhost:5174/posts", {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to fetch posts");
            }

            posts = await response.json();
        } catch (error) {
            console.error("Error fetching posts:", error);
        }
    }

    async function toggleLike(post) {
        try {
            const response = await fetch(`http://localhost:5174/like/${post.id}`, {
                method: "POST",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to toggle like");
            }

            const updatedPost = await response.json();
            const index = posts.findIndex(p => p.id === post.id);
            posts[index] = updatedPost;
            posts = [...posts]; // Trigger reactivity
        } catch (error) {
            console.error("Error toggling like:", error);
        }
    }
</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <h1>Recent Trades</h1>

        <table>
            <thead>
                <tr>
                    <th>Username</th>
                    <th>Ticker</th>
                    <th>Quantity</th>
                    <th>Trade Type</th>
                    <th>Rationale</th>
                    <th>Date</th>
                    <th>Likes</th>
                </tr>
            </thead>
            <tbody>
                {#each posts as post}
                    <tr>
                        <td>{post.username}</td>
                        <td>{post.symbol}</td>
                        <td>{post.quantity}</td>
                        <td>{post.trade_type}</td>
                        <td>{post.rationale}</td>
                        <td>{new Date(post.trade_date).toLocaleString()}</td>
                        <td>
                            <button
                                class="like-heart"
                                on:click={() => toggleLike(post)}
                            >
                                {#if post.liked_by_user}
                                    ♥
                                {:else}
                                    ♡
                                {/if}
                            </button>
                            <span class="like-count">{post.likes}</span>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
    <Footer />
</main>

<style>
    .content {
        margin-top: 96px;
        padding: 20px;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        margin-top: 20px;
    }

    th, td {
        padding: 10px;
        text-align: left;
        border: 1px solid #ddd;
    }

    th {
        background-color: #f4f4f4;
        font-weight: bold;
    }

    tr:nth-child(even) {
        background-color: #f9f9f9;
    }

    tr:hover {
        background-color: #f1f1f1;
    }

    .like-heart {
        cursor: pointer;
        font-size: 24px;
        color: #ff6f61;
        background: none;
        border: none;
        padding: 0;
        transition: color 0.3s;
    }

    .like-heart:hover {
        color: #e63946;
    }

    .like-count {
        font-size: 16px;
        margin-left: 10px;
        color: #333;
    }
</style>