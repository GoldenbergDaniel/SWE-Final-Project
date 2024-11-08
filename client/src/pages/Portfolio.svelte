<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';

    // Example data for the portfolio (this would usually come from an API or data store)
    let stocks = [
        { ticker: 'AAPL', price: 175.43, quantity: 50, date: '2023-08-01' },
        { ticker: 'GOOGL', price: 2840.12, quantity: 20, date: '2023-07-15' },
        { ticker: 'AMZN', price: 3341.59, quantity: 10, date: '2023-06-20' },
        { ticker: 'TSLA', price: 755.87, quantity: 15, date: '2023-05-10' },
        { ticker: 'MSFT', price: 297.95, quantity: 30, date: '2023-04-12' }
    ];
</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <h1>Portfolio Page</h1>

        <!-- Table for displaying stock data -->
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
