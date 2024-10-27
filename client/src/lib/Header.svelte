<script lang=ts>
    import { Link } from 'svelte-routing';
    
    // State to manage hamburger menu and dropdown
    let isMenuOpen = false;
    let isDropdownOpen = false;

    // Toggles the hamburger menu open/close
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

    // Toggles the dropdown open/close
    function toggleDropdown() {
        isDropdownOpen = !isDropdownOpen;
    }

    // Close the dropdown when clicking outside
    function closeDropdown(event) {
        if (!event.target.closest('.dropbtn')) {
            isDropdownOpen = false;
        }
    }

    // Close the menu on window resize (optional, as you were doing)
    function handleResize() {
        if (window.innerWidth >= 1024) {
            isMenuOpen = false;
        }
    }

    window.addEventListener('resize', handleResize);
</script>

<main>
    <div class="header-container">
        <div class="subcontainer">
            <nav class="navbar">
                <Link to="/dashboard" class="logo-container">
                    <img src="tradex_logo.jpg" alt="TradEx Logo">
                </Link>
            </nav>
            <ul class="nav-menu {isMenuOpen ? 'active' : ''}">
                <li class="nav-item">
                    <Link to="/dashboard" class="nav-link">Dashboard</Link>
                </li>
                <li class="nav-item">
                    <Link to="/portfolio" class="nav-link">Portfolio</Link>
                </li>
                <li class="nav-item">
                    <Link to="/posts" class="nav-link">Posts</Link>
                </li>
                <li class="nav-item spaced-nav-item">
                    <Link to="/connections" class="nav-link">Connections</Link>
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

    /* Spacing from the right side of the screen */
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

    @media (max-width: 800px) {
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
            height: 100%; /* Adjusted to take the full screen height */
            text-align: center;
            transition: right 0.3s;
            justify-content: flex-start;
            overflow-y: auto; /* Enables vertical scrolling */
        }

        .nav-item {
            margin: 16px 0;
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