<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { checkAuth } from '../auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    let leaderboardData = [];

    onMount(async () => {
        if (!checkAuth()) {
            alert('Access Denied! Login Required');
            navigate('/');
        } else {
            await fetchLeaderboardData();
        }
    });

    async function fetchLeaderboardData() {
        try {
            const response = await fetch("http://localhost:5174/leaderboard", {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to fetch leaderboard data");
            }

            leaderboardData = await response.json();
        } catch (error) {
            console.error("Error fetching leaderboard data:", error);
        }
    }
</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <div class="dashboard">
            <div class="leaderboard">
                <h2>Leaderboard</h2>
                <table>
                    <thead>
                        <tr>
                            <th>Rank</th>
                            <th>Username</th>
                            <th>Portfolio Value</th>
                            <th>% Gain/Loss</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each leaderboardData as { username, portfolioValue, gainLoss }, index}
                            <tr>
                                <td>{index + 1}</td>
                                <td>{username}</td>
                                <td>${portfolioValue.toFixed(2)}</td>
                                <td class={gainLoss >= 0 ? 'positive' : 'negative'}>
                                    {gainLoss.toFixed(2)}%
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
    <Footer />
</main>

<style>
    .content {
        margin-top: 96px;
        padding: 20px;
    }

    /* Flexbox layout for dashboard */
    .dashboard {
        display: flex;
        justify-content: space-between;
        gap: 20px;
    }

    /* Leaderboard */
    .leaderboard {
        flex: 1;
        min-width: 45%;
        background-color: rgba(59, 47, 47, 0.87);
        border-radius: 8px;
        padding: 20px;
    }

    .leaderboard h2 {
        text-align: center;
        margin-bottom: 20px;
        color: rgb(229, 228, 217);
        font-size: 30px;
    }

    .leaderboard table {
        width: 100%;
        border-collapse: collapse;
    }

    .leaderboard th, .leaderboard td {
        padding: 10px;
        text-align: left;
        border: 1px solid rgb(229, 228, 217);
        color: rgb(229, 228, 217);
        font-size: 23px;
    }

    .leaderboard th {
        background-color: rgb(229, 228, 217);
        color: rgba(59, 47, 47, 0.87);
        font-weight: bold;
    }

    .leaderboard tr:nth-child(even) {
        background-color: rgba(229, 228, 217, 0.1);
    }

    .leaderboard tr:hover {
        background-color: rgba(229, 228, 217, 0.2);
    }
    
    .positive {
        color: #4CAF50 !important;
    }

    .negative {
        color: #FF5252 !important;
    }
</style>