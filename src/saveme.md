import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import Web3 from 'web3';
import './style.css';


const BridgeCrypto = () => {
    const web3 = new Web3();
    const [formData, setFormData] = useState({
        currency: 'octa',
        fromChain: 'octa',
        amount: '1',
        bridgeTo: 'grams',
        shippingAddress: '',
    });

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };
    const [ws, setWs] = useState(null);
    const URL = "ws://192.168.50.90:8080/ws";


    React.useEffect(() => {
        connectMetaMask();
    }, []);


    const initWebSocketConnection = async (url) => {
        const websocket = new WebSocket(url);

        websocket.onopen = () => {
            console.log('WebSocket connection opened');
        };

        websocket.onmessage = (event) => {
            console.log('WebSocket message received:', event.data);
            alert(event.data)
        };

        websocket.onclose = () => {
            console.log('WebSocket connection closed');
        };

        websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        setWs(websocket);
    };


    const connectMetaMask = async () => {
        if (window.ethereum) {
            try {
                const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
                setFormData({ ...formData, shippingAddress: accounts[0] });
                // modify the URL to include the metamask address
                const url = `${URL}?ethaddress=${accounts[0]}`;
                    // initate the websocket connection
                    if (ws === null){
                        initWebSocketConnection(url);
                    }
            } catch (error) {
                console.error('Error:', error);
            }
        } else {
            alert('MetaMask is not installed. Please install MetaMask and try again.');
        }
    };

    const requestChangeToOctaSpaceNetwork = async () => {
        await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: web3.utils.toHex(800001) }],
        });
    };

    const requestChangeToPartyChainNetwork = async () => {
        await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: web3.utils.toHex(1773) }],
        });
    };

    const handleSubmit = async (e) => {
        console.log(formData);
        e.preventDefault();
    
        // Validate that the user has connected MetaMask
        if (!formData.shippingAddress) {
            alert('Please connect MetaMask to continue.')
            connectMetaMask();
            return;
        }
    
        // Validate that the user has selected the correct network
        try {
            if (formData.fromChain === 'octa') {
                await requestChangeToOctaSpaceNetwork();
            }
            if (formData.fromChain === 'grams') {
                await requestChangeToPartyChainNetwork();
            }
        } catch (error) {
            console.error('Error switching networks:', error);
            return;
        }
    
        // Validate that all fields are populated
        for (const field in formData) {
            if (!formData[field]) {
                alert(`Please fill in the ${field} field.`);
                return;
            }
        }
    
        // Convert the amount into a BigInt
        const amountInWei = parseInt(web3.utils.toWei(formData.amount, 'ether'));
    
        console.log("im herrreee");
        try {
            const response = await axios.post(
                '/requestbridge',
                {
                    ...formData,
                    amount: amountInWei,
                },
                {
                    headers: {
                        'Content-Type': 'application/json',
                    },
                },
            ).catch((error) => {
                if (error.response) {
                    console.error("Server responded with an error:", error.response.data);
                } else if (error.request) {
                    console.error("No response was received from the server:", error.request);
                } else {
                    console.error("Axios error:", error.message);
                }
            });
    
            if (!response) {
                console.error("No response object");
                return;
            }
            console.log("im herrreee");
            console.log(response.data)
    
            // Verify the response.data has valid fields
            if (!response.data.address || !response.data.amount) {
                alert('Invalid response from server. Please try again.');
                return;
            }
            console.log("im herrreee");
    
            // Convert the response amount to a hexadecimal string
            var responseAmountHex = web3.utils.toHex(response.data.amount);
    
            // Create a transaction
            try {
                const transaction = await window.ethereum.request({
                    method: 'eth_sendTransaction',
                    params: [
                        {
                            from: formData.shippingAddress,
                            to: response.data.address,
                            value: responseAmountHex,
                        },
                    ],
                });
                console.log(transaction);
                console.log("im herrreee");
            } catch (error) {
                console.error('Error sending transaction:', error);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };
    

    return (
        <div>
            <header>
                <nav class="navbar">
                    <div class="container mx-auto pl-2 pr-2">
                        <div class="grid grid-cols-12 gap-4 content-center">
                            <div class=" col-span-4 navbar__logo">
                                Bridge
                            </div>
                            <div class="col-span-8 flex flex-row justify-end items-center">
                                <button class="btn btn--outline btn--metamask mr-5"
                                    onClick={connectMetaMask}
                                >Connect Metamask</button>
                                <div class="theme-toggle">
                                    <input type="checkbox" class="theme-toggle__input" id="theme-toggle__input" />
                                    <label for="theme-toggle__input" class="theme-toggle__label">
                                        <span class="theme-toggle__icon theme-toggle__dark">
                                            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
                                                <path d="M17.4776 10.2333C17.4172 10.2068 17.3493 10.2024 17.286 10.2211C17.2227 10.2398 17.1681 10.2802 17.1317 10.3352C16.4611 11.3304 15.5539 12.1434 14.4915 12.7015C13.4292 13.2596 12.2449 13.5452 11.045 13.5327C9.10715 13.5167 7.29498 12.7574 5.94224 11.3947C5.26518 10.7137 4.73203 9.90334 4.37453 9.012C4.01702 8.12067 3.84246 7.16652 3.86123 6.20635C3.88 5.24617 4.09172 4.29958 4.48379 3.42289C4.87586 2.54621 5.44027 1.75735 6.14344 1.10325C6.19154 1.05844 6.22245 0.998216 6.23082 0.933013C6.23919 0.86781 6.22449 0.801734 6.18927 0.746229C6.15405 0.690724 6.10052 0.649286 6.03796 0.629095C5.9754 0.608903 5.90775 0.611229 5.84672 0.635669C4.23625 1.27439 2.85241 2.37815 1.87157 3.80627C0.72077 5.48216 0.196743 7.50971 0.391396 9.53333C0.586048 11.557 1.48695 13.4474 2.9361 14.8732C4.58187 16.4929 6.76975 17.3848 9.0966 17.3848C11.1064 17.3843 13.0555 16.6957 14.6196 15.4336C16.1501 14.1949 17.2155 12.474 17.6418 10.5517C17.6563 10.4878 17.6479 10.4208 17.6178 10.3625C17.5878 10.3043 17.5381 10.2585 17.4776 10.2333ZM14.354 15.1059C12.8651 16.3072 11.0097 16.9625 9.0966 16.9629C6.88112 16.9629 4.79836 16.114 3.23201 14.5725C1.85354 13.2163 0.996548 11.4181 0.81134 9.49316C0.626131 7.56827 1.12453 5.63961 2.21912 4.04544C3.01267 2.88994 4.08594 1.95424 5.33882 1.32561C4.06467 2.76954 3.38717 4.64389 3.44362 6.56877C3.50007 8.49366 4.28627 10.3251 5.64285 11.6919C7.07445 13.134 8.99173 13.9376 11.0414 13.9545C11.0631 13.9545 11.0845 13.9548 11.1062 13.9548C12.2414 13.9571 13.3629 13.7074 14.3898 13.2236C15.4167 12.7399 16.3234 12.0341 17.0445 11.1574C16.5538 12.7087 15.6183 14.0817 14.354 15.1058V15.1059Z" fill="white" />
                                            </svg>
                                        </span>
                                        <span class="theme-toggle__icon theme-toggle__light">
                                            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
                                                <g clip-path="url(#clip0_119_125)">
                                                    <path d="M8.98996 14.4581C5.97492 14.4581 3.52179 12.0051 3.52179 8.9899C3.52179 5.97486 5.97492 3.52173 8.98996 3.52173C12.0052 3.52173 14.4581 5.97486 14.4581 8.9899C14.4581 12.0051 12.0052 14.4581 8.98996 14.4581ZM8.98996 3.91688C6.19265 3.91688 3.91695 6.19259 3.91695 8.9899C3.91695 11.787 6.19265 14.0629 8.98996 14.0629C11.7871 14.0629 14.063 11.787 14.063 8.9899C14.063 6.19259 11.7871 3.91688 8.98996 3.91688Z" fill="black" />
                                                    <path d="M8.98994 2.28271C8.88087 2.28271 8.79236 2.19419 8.79236 2.08513V0.296455C8.79236 0.187392 8.88087 0.098877 8.98994 0.098877C9.099 0.098877 9.18751 0.187392 9.18751 0.296455V2.08513C9.18751 2.19419 9.099 2.28271 8.98994 2.28271Z" fill="black" />
                                                    <path d="M8.98994 17.8809C8.88087 17.8809 8.79236 17.7925 8.79236 17.6833V15.8946C8.79236 15.7853 8.88087 15.697 8.98994 15.697C9.099 15.697 9.18751 15.7853 9.18751 15.8946V17.6833C9.18751 17.7925 9.099 17.8809 8.98994 17.8809Z" fill="black" />
                                                    <path d="M2.08519 9.18739H0.296516C0.187453 9.18739 0.098938 9.09888 0.098938 8.98981C0.098938 8.88075 0.187453 8.79224 0.296516 8.79224H2.08519C2.19425 8.79224 2.28277 8.88075 2.28277 8.98981C2.28277 9.09888 2.19425 9.18739 2.08519 9.18739Z" fill="black" />
                                                    <path d="M17.6834 9.18739H15.8947C15.7855 9.18739 15.6971 9.09888 15.6971 8.98981C15.6971 8.88075 15.7855 8.79224 15.8947 8.79224H17.6834C17.7927 8.79224 17.881 8.88075 17.881 8.98981C17.881 9.09888 17.7927 9.18739 17.6834 9.18739Z" fill="black" />
                                                    <path d="M13.8723 4.30517C13.8217 4.30517 13.7712 4.2858 13.7326 4.24728C13.6554 4.17002 13.6554 4.04515 13.7326 3.9679L14.9973 2.7032C15.0746 2.62595 15.1994 2.62595 15.2767 2.7032C15.3539 2.78046 15.3539 2.90533 15.2767 2.98258L14.012 4.24728C13.9735 4.2858 13.9229 4.30517 13.8723 4.30517Z" fill="black" />
                                                    <path d="M2.84271 15.3345C2.79213 15.3345 2.74155 15.3151 2.70302 15.2766C2.62577 15.1993 2.62577 15.0745 2.70302 14.9972L3.96772 13.7325C4.04497 13.6552 4.16984 13.6552 4.24709 13.7325C4.32435 13.8098 4.32435 13.9346 4.24709 14.0119L2.9824 15.2766C2.94387 15.3153 2.89329 15.3345 2.84271 15.3345Z" fill="black" />
                                                    <path d="M4.10759 4.30517C4.05701 4.30517 4.00643 4.2858 3.9679 4.24728L2.7032 2.98258C2.62595 2.90533 2.62595 2.78046 2.7032 2.7032C2.78046 2.62595 2.90533 2.62595 2.98258 2.7032L4.24728 3.9679C4.32453 4.04515 4.32453 4.17002 4.24728 4.24728C4.20855 4.2858 4.15797 4.30517 4.10759 4.30517Z" fill="black" />
                                                    <path d="M15.137 15.3345C15.0864 15.3345 15.0358 15.3151 14.9973 15.2766L13.7326 14.0119C13.6554 13.9346 13.6554 13.8098 13.7326 13.7325C13.8099 13.6552 13.9347 13.6552 14.012 13.7325L15.2767 14.9972C15.3539 15.0745 15.3539 15.1993 15.2767 15.2766C15.2382 15.3153 15.1876 15.3345 15.137 15.3345Z" fill="black" />
                                                </g>
                                                <defs>
                                                    <clipPath id="clip0_119_125">
                                                        <rect width="18" height="18" fill="white" />
                                                    </clipPath>
                                                </defs>
                                            </svg>
                                        </span>
                                        <span class="theme-toggle__ball"></span>
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>
                </nav>
            </header>
            <div className="container mx-auto pl-2 pr-2">
                <div className="sm:grid grid-cols-12 gap-5 content-center">

                    <div className="xl:col-span-4 xl:col-start-2 sm:col-span-5">
                        <div className="box relative z-30">

                            <div className="box__select-group">
                                <label htmlFor="choose-network-1">Choose network</label>
                                <select name="fromChain" id="fromChain"
                                    onChange={handleChange}>
                                    <option value="octa">OCTA NETWORK</option>
                                    <option value="grams">PARTYCHAIN</option>
                                </select>
                            </div>

                            <div className="box__select-group">
                                <label htmlFor="choose-coin-1">Choose coin</label>
                                <select name="currency" id="currency"
                                    onChange={handleChange}
                                >
                                    <option value="octa">OCTA</option>
                                    <option value="grams">GRAM</option>
                                    <option value="wocta">wOCTA</option>
                                    <option value="wgrams">wGRAM</option>
                                </select>
                            </div>

                            <div className="box__input-group">
                                <label htmlFor="amount" className="mr-2">Amount:</label>
                                <input name="amount" type="number" id="amount" placeholder="100 ... 100000" min="100" max="100000"
                                    value={formData.amount}
                                    onChange={handleChange}
                                />
                            </div>

                        </div>
                        <div className="box box--small">
                            <p className="text-lg">Balance: <span className="font-semibold"></span></p>
                            <p className="text-lg">Fee: <span className="font-semibold"></span></p>
                            <p className="text-lg">Minimum: <span className="font-semibold"></span></p>
                        </div>
                    </div>

                    <div className="sm:col-span-2 sm:text-center">
                        <div className="info pt-12 sm:pb-0 pb-12 sm:flex flex-col h-full">
                            <p className="text-2xl font-semibold  sm:pl-0 pl-20">
                                Swap direction
                            </p>
                            <div className="arrow relative sm:block hidden">
                                <svg className="z-20 arrow__desktop" width="292" height="36" viewBox="0 0 292 36" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M292 18L262 0.679492V35.3205L292 18ZM0 21H265V15H0L0 21Z" fill="#BE6432" />
                                </svg>
                            </div>
                            <p className="pt-12 mt-10 font-bold text-2xl sm:pl-0 pl-20">Monitor</p>
                            <p className="mt-2 text-lg sm:pl-0 pl-20 mb-5 sm:mb-0">Loading...</p>
                            <button className="btn btn--main sm:ml-0 ml-20"
                                onClick={handleSubmit}
                            >Request swap</button>
                        </div>

                    </div>

                    <div className="xl:col-span-4 sm:col-span-5">
                        <div className="box relative z-10">

                            <div className="box__select-group">
                                <label htmlFor="bridgeTo">Choose network</label>
                                <select name="bridgeTo" id="bridgeTo"
                                    onChange={handleChange}>
                                    <option value="grams">PARTYCHAIN</option>
                                    <option value="octa">OCTA NETWORK</option>
                                </select>
                            </div>

                            {/* <div className="box__select-group">
                                <label htmlFor="choose-coin-2">Choose coin</label>
                                <select name="choose-coin-2" id="choose-coin-2">
                                    <option value="gram">GRAM</option>
                                    <option value="octa">OCTA</option>
                                </select>
                            </div> */}

                            {/* if the selected network is grams show wGRAMS if it is octa then choose wOCTA */}
                            {/* <p className="text-xl pt-4">
                                Will be received: <span className="font-bold text-2xl text-accent-primary">3000</span>{" "}
                                {formData.fromChain === "grams" ? "wGRAMS" : formData.fromChain === "octa" ? "wOCTA" : ""}
                            </p> */}

                        </div>
                        <div className="box box--small flex items-center">
                            <p className="text-lg">Balance: <span className="font-semibold"></span></p>
                        </div>
                    </div>

                </div>
            </div>
        </div >
    );
};

export default BridgeCrypto;