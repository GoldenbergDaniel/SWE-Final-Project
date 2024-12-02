<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { checkAuth } from '../auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    let userBalance = 0;
    let ticker = '';
    let shareNumber = 0;
    let tradeType = 'buy';
    let rationale = '';

    onMount(async () => {
        if (!checkAuth()) {
            alert('Access Denied! Login Required');
            navigate('/');
        } else {
            await fetchUserData();
        }
    });

    async function fetchUserData() {
        try {
            const response = await fetch("http://localhost:5174/userdata", {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to fetch user data");
            }

            const data = await response.json();
            userBalance = data.balance;
        } catch (error) {
            console.error("Error fetching user data: ", error);
        }
    }

    function handleShareNumberInput(event) {
        if (shareNumber < 0) {
            shareNumber = 0;
        }
    }

    async function placeOrder() {
        try {
            const response = await fetch("http://localhost:5174/trade", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    symbol: ticker,
                    quantity: shareNumber,
                    trade_type: tradeType,
                    rationale: rationale
                }),
                credentials: "include",
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || "Failed to place order");
            }

            const data = await response.json();
            alert(data.message);

            // Update user balance
            userBalance = data.new_balance;

            // Reset form fields
            ticker = '';
            shareNumber = 0;
            rationale = '';

        } catch (error) {
            console.error("Error placing order:", error);
            alert(error.message); }
    }
</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <h1>Trade Page</h1>
    </div>
    <div class="trade">
        <h2>Your Balance</h2>
        <table>
            <thead>
                <tr>
                    <th>Balance</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>${userBalance.toFixed(2)}</td>
                </tr>
            </tbody>
        </table>

        <h2>Make a Trade</h2>
        <form>
            <input
                type="text"
                placeholder="TICKER"
                bind:value={ticker}
                class="input-field"
            />
            <input
                type="number"
                placeholder="Share Number"
                bind:value={shareNumber}
                min="0"
                on:input={handleShareNumberInput}
                class="input-field"
            />
            <select bind:value={tradeType} class="input-field">
                <option value="buy">Buy</option>
                <option value="sell">Sell</option>
            </select>
            <input
                type="text"
                placeholder="Rationale"
                bind:value={rationale}
                class="input-field"
            />
            <button type="button" on:click={placeOrder} class="place-order-btn">Place Order</button>
        </form>
    </div>
    <Footer />
</main>

<style>
    .content{
        margin-top: 96px;
    }
    /* Right side: Balance and Trade form */
    .trade {
        flex: 0 0 35%;
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
        background-color: #4CAF50;
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
    select.input-field {
        appearance: none;
        background-image: url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23007CB2%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E");
        background-repeat: no-repeat;
        background-position: right 0.7em top 50%;
        background-size: 0.65em auto;
    }
</style>