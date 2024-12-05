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
            const errorText = await response.text();
            console.error("Error response:", errorText);
            throw new Error(`Failed to place order: ${response.status} ${response.statusText}`);
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
        alert(error.message);
    }
}
</script>

<main>
    <Background />
    <Header />
    <div class="trade">
        <h2>Your Balance</h2>
        <table class="balance-table">
            <thead>
                <tr>
                    <th>Cash Available for Trading</th>
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
    input {
        color: rgba(59, 47, 47, 0.87);
        background-color:rgb(229, 228, 217);
    }

    .trade {
        flex: 0 0 35%;
        min-width: 250px;
        padding: 20px;
        border: 1px solid #ddd;
        background-color: rgba(59, 47, 47, 0.87);
        border-radius: 8px;
        margin-top: 110px;
    }

    .trade h2 {
        text-align: center;
        margin-bottom: 20px;
        color: rgb(229, 228, 217);
    }

    .trade table {
        width: 100%;
        border-collapse: collapse;
        margin-bottom: 20px;
    }

    .trade table thead {
        color: rgb(229, 228, 217);
        background-color: rgb(229, 228, 217);
    }

    .trade table th {
        background-color: rgb(229, 228, 217);
        color: rgba(59, 47, 47, 0.87);
        padding: 10px;
        text-align: center;
        border: 1px solid #ddd;
    }

    .trade th {
        background-color: #f4f4f4;
        font-weight: bold;
    }

    .trade table td {
        padding: 10px;
        text-align: center;
        border: 1px solid #ddd;
        color: rgb(229, 228, 217);
        font-size: 1.2em;
    }

    .trade td {
        font-size: 1.2em;
    }

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
    select{
        background-color: rgb(229, 228, 217);
        color: rgba(59, 47, 47, 0.87);
    }
    select.input-field {
        appearance: none;
        background-image: url('data:image/svg+xml;charset=US-ASCII,%3Csvg%20fill%3D%22rgba%2859%2C%2047%2C%2047%2C%200.87%29%22%20viewBox%3D%22-6.5%200%2032%2032%22%20version%3D%221.1%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20stroke%3D%22rgba%2859%2C%2047%2C%2047%2C%200.87%29%22%3E%3Cg%20id%3D%22SVGRepo_bgCarrier%22%20stroke-width%3D%220%22%3E%3C%2Fg%3E%3Cg%20id%3D%22SVGRepo_tracerCarrier%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%3E%3C%2Fg%3E%3Cg%20id%3D%22SVGRepo_iconCarrier%22%3E%3Ctitle%3Edropdown%3C%2Ftitle%3E%3Cpath%20d%3D%22M18.813%2011.406l-7.906%209.906c-0.75%200.906-1.906%200.906-2.625%200l-7.906-9.906c-0.75-0.938-0.375-1.656%200.781-1.656h16.875c1.188%200%201.531%200.719%200.781%201.656z%22%3E%3C%2Fpath%3E%3C%2Fg%3E%3C%2Fsvg%3E');
        background-repeat: no-repeat;
        background-position: right 0.7em top 50%;
        background-size: 1.20em auto;
    }
</style>