import { hydrateRoot } from 'react-dom/client';
import React, { Component } from 'react';
document.body.innerHTML = '<div id="app"></div>';
import { Header } from './src/header';
import { Info } from './src/info';
import { Footer } from './src/footer';

const root = hydrateRoot(document.getElementById('app'),
	<div className="h-screen w-screen">
		<Header />
		<Info />
		<Footer />
	</div>
);
