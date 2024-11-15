<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { checkAuth} from '../assets/auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    onMount(() => {
        if (!checkAuth()) {
            alert('Access Denied! Login Required');
            navigate('/');
        }
    });


    // Example data for the leaderboard
    let tableData = [
        { username: 'user1', portfolioValue: 12000, gainLoss: 5 },
        { username: 'user2', portfolioValue: 9000, gainLoss: -3 },
        { username: 'user3', portfolioValue: 15000, gainLoss: 8 },
        { username: 'user4', portfolioValue: 7000, gainLoss: -2 },
        { username: 'user5', portfolioValue: 12000, gainLoss: 4 },
        { username: 'user6', portfolioValue: 8000, gainLoss: 6 },
        { username: 'user7', portfolioValue: 10000, gainLoss: 2 },
        { username: 'user8', portfolioValue: 11000, gainLoss: -1 },
        { username: 'user9', portfolioValue: 13000, gainLoss: 7 },
        { username: 'user10', portfolioValue: 14000, gainLoss: 3 }
    ];

</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <!-- Flexbox container for layout -->
        <div class="dashboard">
            <!-- Left side: Leaderboard -->
            <div class="leaderboard">
                <h2>Leaderboard</h2>
                <table>
                    <thead>
                        <tr>
                            <th>Username</th>
                            <th>Portfolio Value</th>
                            <th>% Gain/Loss</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each tableData as { username, portfolioValue, gainLoss }}
                            <tr>
                                <td>{username}</td>
                                <td>${portfolioValue}</td>
                                <td>{gainLoss}%</td>
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
        gap: 20px; /* Space between the leaderboard and trade form */
    }

    /* Left side: Leaderboard */
    .leaderboard {
        flex: 1;
        min-width: 45%; /* Ensure the leaderboard takes up less than half of the screen */
    }

    .leaderboard h2 {
        text-align: center;
        margin-bottom: 20px;
    }

    .leaderboard table {
        width: 100%;
        border-collapse: collapse;
    }

    .leaderboard th, .leaderboard td {
        padding: 10px;
        text-align: left;
        border: 1px solid #ddd;
    }

    .leaderboard th {
        background-color: #f4f4f4;
        font-weight: bold;
    }

    .leaderboard tr:nth-child(even) {
        background-color: #f9f9f9;
    }

    .leaderboard tr:hover {
        background-color: #f1f1f1;
    }

    /* Right side: Balance and Trade form */
    .trade {
        flex: 0 0 35%; /* Trade section takes less space */
        min-width: 250px;
        padding: 20px;
        border: 1px solid #ddd;
        background-color: #f9f9f9;
        border-radius: 8px;
    }

    .trade h2 {
        text-align: center;
        margin-bottom: 20px;
    }

    /* Balance table */
    .trade table {
        width: 100%;
        border-collapse: collapse;
        margin-bottom: 20px;
    }

    .trade th, .trade td {
        padding: 10px;
        text-align: center;
        border: 1px solid #ddd;
    }

    .trade th {
        background-color: #f4f4f4;
        font-weight: bold;
    }

    .trade td {
        font-size: 1.2em;
    }

    /* Input fields */
    .input-field {
        width: 100%;
        padding: 10px;
        margin-bottom: 15px;
        border: 1px solid #ddd;
        border-radius: 5px;
        font-size: 16px;
        box-sizing: border-box;
    }

    .input-field::placeholder {
        color: #aaa;
    }

    .input-field:focus {
        border-color: #4CAF50;
        outline: none;
    }

    /* Place Order button */
    .place-order-btn {
        width: 100%;
        padding: 12px;
        background-color: #4CAF50; /* Green */
        color: white;
        border: none;
        border-radius: 5px;
        font-size: 16px;
        cursor: pointer;
        transition: background-color 0.3s;
    }

    .place-order-btn:hover {
        background-color: #45a049;
    }

    .place-order-btn:active {
        background-color: #3e8e41;
    }
</style>
