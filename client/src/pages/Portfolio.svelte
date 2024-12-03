<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { checkAuth } from '../auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    let portfolioData = null;
    let isLoading = true;

    onMount(async () => {
        if (!checkAuth()) {
            alert('Access Denied! Login Required');
            navigate('/');
        } else {
            await fetchPortfolioData();
        }
    });

    async function fetchPortfolioData() {
        try {
            const response = await fetch("http://localhost:5174/portfolio-value", {
                method: "GET",
                credentials: "include",
            });

            console.log("Response status:", response.status);
            console.log("Response headers:", response.headers);

            if (!response.ok) {
                const errorText = await response.text();
                console.error("Error response:", errorText);
                throw new Error(`Failed to fetch portfolio data: ${response.status} ${response.statusText}`);
            }

            portfolioData = await response.json();
            console.log("Received portfolio data:", portfolioData);
        } catch (error) {
            console.error("Error fetching portfolio data:", error);
        } finally {
            isLoading = false;
        }
    }

    function formatCurrency(value) {
        return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value);
    }
</script>

<main>
    <Background />
    <Header />
    <div class="content">
        {#if isLoading}
            <p>Loading portfolio data...</p>
        {:else if portfolioData}
            <div class="tables">
                <!-- Stock data table -->
                <div class="table-container">
                    <h2>Stock Portfolio</h2>
                    <table>
                        <thead>
                            <tr>
                                <th>Ticker</th>
                                <th>Current Price</th>
                                <th>Quantity</th>
                                <th>Value</th>
                                <th>Avg. Price</th>
                                <th>Profit/Loss</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each Object.entries(portfolioData.portfolio) as [ticker, stock]}
                                <tr>
                                    <td>{ticker}</td>
                                    <td>{formatCurrency(stock.currentPrice)}</td>
                                    <td>{stock.quantity}</td>
                                    <td>{formatCurrency(stock.currentValue)}</td>
                                    <td>{formatCurrency(stock.averagePrice)}</td>
                                    <td class={stock.profitLoss >= 0 ? 'profit' : 'loss'}>
                                        {formatCurrency(stock.profitLoss)}
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                    <p><strong>Total Portfolio Value: {formatCurrency(portfolioData.totalValue)}</strong></p>
                </div>

                <!-- User data table -->
                <div class="table-container">
                    <h2>User Information</h2>
                    <table>
                        <thead>
                            <tr>
                                <th>Field</th>
                                <th>Value</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td>Username</td>
                                <td>{portfolioData.username}</td>
                            </tr>
                            <tr>
                                <td>Email</td>
                                <td>{portfolioData.email}</td>
                            </tr>
                            <tr>
                                <td>Cash Balance</td>
                                <td>{formatCurrency(portfolioData.balance)}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        {:else}
            <p>Failed to load portfolio data. Please try again later.</p>
        {/if}
    </div>
    <Footer />
</main>

<style>
    .content {
        margin-top: 96px;
        padding: 20px;
    }

    .tables {
        display: flex;
        justify-content: space-between;
        gap: 60px;
    }

    .table-container {
        width: 60%;
        background-color: rgba(59, 47, 47, 0.87);
        border-radius: 8px;
        padding: 20px;
    }

    h1, h2 {
        color: rgb(229, 228, 217);
        text-align: center;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        margin-top: 20px;
    }

    th, td {
        padding: 10px;
        text-align: left;
        border: 1px solid rgb(229, 228, 217);
        color: rgb(229, 228, 217);
    }

    th {
        background-color: rgb(229, 228, 217);
        color: rgba(59, 47, 47, 0.87);
        font-weight: bold;
    }

    tr:nth-child(even) {
        background-color: rgba(229, 228, 217, 0.1);
    }

    tr:hover {
        background-color: rgba(229, 228, 217, 0.2);
    }

    .profit {
        color: #4CAF50;
    }

    .loss {
        color: #FF5252;
    }

    p {
        color: rgb(229, 228, 217);
        text-align: center;
        font-weight: bold;
        margin-top: 20px;
    }
</style>