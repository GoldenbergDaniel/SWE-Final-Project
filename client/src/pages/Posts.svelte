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
            console.log("Fetched posts:", posts);
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

            const updatedLikeData = await response.json();
            const index = posts.findIndex(p => p.id === post.id);
            posts[index] = {
                ...posts[index],
                likes_count: updatedLikeData.likes_count,
                liked_by_user: updatedLikeData.liked_by_user
            };
            posts = [...posts]; // Trigger reactivity
        } catch (error) {
            console.error("Error toggling like:", error);
        }
    }
</script>

<main>
    <Background />
    <Header />
    <div class="posts-container">
        <h2>Recent Trades</h2>
        {#if posts.length === 0}
            <p class="no-posts">No trades have been made yet.</p>
        {:else}
            {#each posts as post}
                <div class="post">
                    <p class="trade-info">
                        <strong>{post.username}</strong> {post.trade_type === 'buy' ? 'bought' : 'sold'} {post.quantity} shares of {post.symbol}
                    </p>
                    <p class="rationale">Rationale: {post.rationale}</p>
                    <p class="date">Date: {new Date(post.trade_date).toLocaleString()}</p>
                    <div class="like-section">
                        <button
                            class="like-button"
                            on:click={() => toggleLike(post)}
                        >
                            {post.liked_by_user ? '♥' : '♡'}
                        </button>
                        <span class="like-count">Likes: {post.likes !== undefined ? post.likes : 'N/A'}</span>
                    </div>
                </div>
            {/each}
        {/if}
    </div>
    <Footer />
</main>

<style>
    .posts-container {
        max-width: 800px;
        margin: 110px auto 0;
        padding: 20px;
        background-color: rgba(59, 47, 47, 0.87);
        border-radius: 8px;
    }

    h2 {
        color: rgb(229, 228, 217);
        text-align: center;
        margin-bottom: 20px;
    }

    .post {
        background-color: rgb(229, 228, 217);
        color: rgba(59, 47, 47, 0.87);
        padding: 15px;
        margin-bottom: 15px;
        border-radius: 5px;
    }

    .post p {
        margin: 5px 0;
    }

    .trade-info {
        font-size: 18px;
    }

    .rationale {
        font-style: italic;
    }

    .date {
        font-size: 14px;
        color: #666;
    }

    .like-section {
        margin-top: 10px;
    }

    .like-button {
        cursor: pointer;
        font-size: 24px;
        color: #ff6f61;
        background: none;
        border: none;
        padding: 0;
        transition: color 0.3s;
    }

    .like-button:hover {
        color: #e63946;
    }

    .like-count {
        font-size: 16px;
        margin-left: 10px;
        color: #333;
    }

    .no-posts {
        color: rgb(229, 228, 217);
        text-align: center;
    }
</style>