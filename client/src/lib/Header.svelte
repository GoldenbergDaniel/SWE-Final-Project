<script lang=ts>
    import { Link } from 'svelte-routing';
    
    let isMenuOpen = false;
    let isDropdownOpen = false;

    function toggleMenu() {
        isMenuOpen = !isMenuOpen;
    }

    function recolorBarsOnHoverOrFocus() {
        const bars = document.querySelectorAll('.bar');
        bars.forEach(bar => {
            (bar as HTMLElement).style.backgroundColor = '#4CAF50'; // Type assertion for HTMLElement
        });
    }

    function resetBarColors() {
        const bars = document.querySelectorAll('.bar');
        bars.forEach(bar => {
            (bar as HTMLElement).style.backgroundColor = 'rgba(59, 47, 47, 0.87)'; // Type assertion for HTMLElement
        });
    }

    function toggleDropdown() {
        isDropdownOpen = !isDropdownOpen;
    }

    function closeDropdown(event) {
        if (!event.target.closest('.dropbtn')) {
            isDropdownOpen = false;
        }
    }

    function handleResize() {
        if (window.innerWidth >= 1024) {
            isMenuOpen = false;
        }
    }

    window.addEventListener('resize', handleResize);

    async function makeLogoutRequest()
    {
        const response = await fetch("http://localhost:5174/logout", {
            credentials: "include",
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: "",
        })

        const responseText = await response.text()
        if (!response.ok)
        {
            throw new Error(responseText || "Logout failed!!!!!!!")
        }

        location.reload()

        return responseText
    }

    function handleLogout()
    {
        makeLogoutRequest()
    }

</script>

<main>
    <div class="header-container">
        <div class="subcontainer">
            <nav class="navbar">
                <Link to="/portfolio" class="logo-container">
                    <img src="tradex_logo.jpg" alt="TradEx Logo">
                </Link>
            </nav>
            <ul class="nav-menu {isMenuOpen ? 'active' : ''}">
                <li class="nav-item">
                    <Link to="/portfolio" class="nav-link">Portfolio</Link>
                </li>
                <li class="nav-item">
                    <Link to="/trade" class="nav-link">Trade</Link>
                </li>
                <li class="nav-item">
                    <Link to="/posts" class="nav-link">Posts</Link>
                </li>
                <li class="nav-item">
                    <Link to="/leaderboard" class="nav-link">Leaderboard</Link>
                </li>
                <li class="nav-item spaced-nav-item">
                    <button type="button" class="nav-link" on:click={handleLogout} aria-label="Logout">Logout</button>
                </li>
            </ul>
            <button 
                type="button" 
                class="hamburger-menu" 
                class:active={isMenuOpen} 
                on:click={toggleMenu} 
                aria-label="Toggle menu"
                on:mouseover={recolorBarsOnHoverOrFocus} 
                on:mouseout={resetBarColors}
                on:focus={recolorBarsOnHoverOrFocus}
                on:blur={resetBarColors}
            >
                <span class="bar"></span>
                <span class="bar"></span>
                <span class="bar"></span>
            </button>
        </div>
    </div>
</main>


<style>
    main {
        position: absolute;
        top: 0;
        left: 0;
        height: 90px;
        width: 100%;
    }

    li {
        list-style: none;
    }

    ul {
        padding-left: 0px;
    }

    img {
        margin-left: 20px;
        height: 90px;
    }

    .subcontainer {
        display: flex;
    }

    .navbar {
        width: 100%;
        min-height: 70px;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .nav-menu {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 60px;
        font-size: 20px;
    }

    .nav-item {
        position: relative;
    }

    .nav-item button {
        color: rgb(229, 228, 217);
        background-color:rgba(59, 47, 47, 0.87);
    }

    .nav-item button:hover{
        color: #c3112c;
    }

    button{
        background-color: rgb(229, 228, 217);
    }

    button:focus{
        outline: none;
        border-color: rgb(229, 228, 217);
    }

    .nav-menu li{
        color: rgba(59, 47, 47, 0.87);
    }

    .spaced-nav-item {
        margin-right: 80px;
    }

    .hamburger-menu {
        display: none;
        cursor: pointer;
    }

    .bar {
        display: block;
        width: 40px;
        height: 3px;
        margin: 5px 0;
        transition: all 0.3s ease;
        background-color: rgba(59, 47, 47, 0.87);
    }

    .bar:hover {
        background-color: #4CAF50;
    }

    @media (max-width: 950px) {
        .hamburger-menu {
            display: block;
            margin-right: 10px;
        }

        .hamburger-menu.active .bar:nth-child(2) {
            opacity: 0;
        }

        .hamburger-menu.active .bar:nth-child(1) {
            transform: translateY(8px) rotate(45deg);
        }

        .hamburger-menu.active .bar:nth-child(3) {
            transform: translateY(-8px) rotate(-45deg);
        }

        .nav-menu {
            position: fixed;
            right: -100%;
            top: 70px;
            gap: 0;
            flex-direction: column;
            background-color: rgba(59, 47, 47, 0.87);
            width: 200px;
            height: 100%;
            text-align: center;
            transition: right 0.3s;
            justify-content: flex-start;
            overflow-y: auto;
        }

        .nav-item {
            margin: 16px 0;
        }

        .nav-item button{
            background-color: transparent;
            border: none;
            cursor: pointer;
            font-size: 20px;
            color: rgb(229, 228, 217);
            padding: 0px;
        }

        .nav-menu.active {
            right: 0;
        }
    }

    @media (max-height: 340px) {
        .nav-menu {
            max-height: 60%;
        }
    }

</style>