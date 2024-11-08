<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';

    // Example data for the table
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
        <h1>Leaderboard</h1>

        <!-- Table for displaying portfolio information -->
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
</style>
