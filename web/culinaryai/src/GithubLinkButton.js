import React from 'react';
import GithubSvg from './github.svg';

const GithubLinkButton = () => (
    <a
        href="https://github.com/windoze95/culinaryai"
        target="_blank"
        className="btn waves-effect waves-light"
        style={{
        marginTop: '16px',
        backgroundColor: 'white',
        color: 'black',
        borderRadius: '12px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center'
        }}
        rel="noopener noreferrer"
    >
        <img
            style={{
                width: '24px',
                height: '24px',
                marginRight: '8px'
            }}
            src={GithubSvg}
            alt="GitHub Logo"
        />
        <span>Follow this on GitHub</span>
    </a>
);

export default GithubLinkButton;