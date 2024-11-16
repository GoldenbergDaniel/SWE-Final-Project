<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { checkAuth } from '../auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    let posts = [
        { username: 'john_doe', ticker: 'AAPL', quantity: 50, rationale: 'Good long-term growth', likes: 0, liked: false },
        { username: 'jane_smith', ticker: 'GOOGL', quantity: 20, rationale: 'Strong market position', likes: 0, liked: false },
        { username: 'alice_jones', ticker: 'AMZN', quantity: 10, rationale: 'Innovative services', likes: 0, liked: false },
        { username: 'bob_miller', ticker: 'TSLA', quantity: 15, rationale: 'Electric vehicles market leader', likes: 0, liked: false },
        { username: 'charlie_brown', ticker: 'MSFT', quantity: 30, rationale: 'Dominance in software industry', likes: 0, liked: false }
    ];

    onMount(() => {
        if (!checkAuth()) {
            alert('Access Denied! Login Required');
            navigate('/');
        }
    });

    const toggleLike = (post) => {
        post.liked = !post.liked;
        if (post.liked)
        {
            post.likes = post.likes + 1;
        }
        else
        {
            post.likes = post.likes - 1;
        }

    };
</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <h1>Posts Page</h1>

        <table>
            <thead>
                <tr>
                    <th>Username</th>
                    <th>Ticker</th>
                    <th>Quantity</th>
                    <th>Rationale</th>
                    <th>Likes</th>
                </tr>
            </thead>
            <tbody>
                {#each posts as post}
                    <tr>
                        <td>{post.username}</td>
                        <td>{post.ticker}</td>
                        <td>{post.quantity}</td>
                        <td>{post.rationale}</td>
                        <td>

                            <button
                                class="like-heart"
                                on:click={() => toggleLike(post)}
                            >
                                {#if post.liked}
                                    ♥
                                {:else}
                                    ♡
                                {/if}
                            </button>
                            <!-- Likes count -->
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

    /* Heart button styles */
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

    /* Rationale column should take up more space */
    td:nth-child(4) {
        width: 40%;
    }
</style>
