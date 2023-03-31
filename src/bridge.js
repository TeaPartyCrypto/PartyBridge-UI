import React, { useState } from 'react';
import axios from 'axios';
import Web3 from 'web3';

const BridgeCrypto = () => {
    const [formData, setFormData] = useState({
        currency: 'grams',
        fromChain: 'grams',
        amount: '1',
        bridgeTo: 'octa',
        shippingAddress: '',
    });

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        // validate that all fields are populated


        // Convert the amount into a BigInt
        const web3 = new Web3();
        const amountInWei = parseInt(web3.utils.toWei(formData.amount, 'ether'))

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
            );
            console.log(response.data);
            alert(JSON.stringify(response.data));
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <div>
            <h1>PartyBridge</h1>
            <form onSubmit={handleSubmit}>
                <label htmlFor="currency">Currency:</label>
                <select
                    id="currency"
                    name="currency"
                    value={formData.currency}
                    onChange={handleChange}
                >
                    <option value="grams">GRAMS</option>
                    <option value="octa">OCTA</option>
                    <option value="wgrams">WGRAMS</option>
                    <option value="wocta">WOCTA</option>
                </select>
                <br />
                <label htmlFor="bridgeTo">Bridge To:</label>
                <select
                    id="bridgeTo"
                    name="bridgeTo"
                    value={formData.bridgeTo}
                    onChange={handleChange}
                >
                    <option value="octa">OctaSpace</option>
                    <option value="grams">PartyChain</option>
                </select>

                <br />
                <label htmlFor="fromChain">From Chain:</label>
                <select
                    id="fromChain"
                    name="fromChain"
                    value={formData.fromChain}
                    onChange={handleChange}
                >
                    <option value="grams">PartyChain</option>
                    <option value="octa">OctaSpace</option>
                </select>
                <br />
                <label htmlFor="amount">Amount:</label>
                <input
                    type="number"
                    id="amount"
                    name="amount"
                    value={formData.amount}
                    onChange={handleChange}
                />

                <br />
                <label htmlFor="shippingAddress">Shipping Address:</label>
                <input
                    type="text"
                    id="shippingAddress"
                    name="shippingAddress"
                    value={formData.shippingAddress}
                    onChange={handleChange}
                />
                <br />
                <button type="submit">Submit</button>
            </form>
        </div>
    );
};

export default BridgeCrypto;

