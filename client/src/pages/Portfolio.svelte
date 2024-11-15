<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { checkAuth } from '../assets/auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    onMount(() => {
        if (!checkAuth()) {
            alert('Access Denied! Login Required');
            navigate('/');
        }
    });

    // Example data for the portfolio (this would usually come from an API or data store)
    let stocks = [
        { ticker: 'AAPL', price: 175.43, quantity: 50, date: '2023-08-01' },
        { ticker: 'GOOGL', price: 2840.12, quantity: 20, date: '2023-07-15' },
        { ticker: 'AMZN', price: 3341.59, quantity: 10, date: '2023-06-20' },
        { ticker: 'TSLA', price: 755.87, quantity: 15, date: '2023-05-10' },
        { ticker: 'MSFT', price: 297.95, quantity: 30, date: '2023-04-12' }
    ];

    // Example user data (replace this with dynamic data later)
    let user = {
        username: 'john_doe',
        email: 'johndoe@example.com'
    };
</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <h1>Portfolio Page</h1>

        <div class="tables">
            <!-- Stock data table -->
            <div class="table-container">
                <h2>Stock Portfolio</h2>
                <table>
                    <thead>
                        <tr>
                            <th>Ticker</th>
                            <th>Price</th>
                            <th>Quantity</th>
                            <th>Date</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each stocks as { ticker, price, quantity, date }}
                            <tr>
                                <td>{ticker}</td>
                                <td>${price.toFixed(2)}</td> <!-- Format the price to two decimal places -->
                                <td>{quantity}</td>
                                <td>{date}</td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
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
                            <td>{user.username}</td>
                        </tr>
                        <tr>
                            <td>Email</td>
                            <td>{user.email}</td>
                        </tr>
                    </tbody>
                </table>
            <button class="logout">Logout</button>
            </div>
        </div>
    </div>
    <Footer />
</main>

<style>

button {
    padding: 0.7em 1.2em;
    font-size: 1em;
    margin: 2em 0;
    background-color: red;
    color: black;
  }

  .button-container {
    display: flex;
    justify-content: space-between;
  }


    .content {
        margin-top: 96px;
        padding: 20px;
    }

    .tables {
        display: flex;
        justify-content: space-between;
        gap: 40px; /* Adjust the gap between tables */
    }

    .table-container {
        width: 48%; /* Each table will take up roughly half the space */
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
