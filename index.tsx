import { hydrateRoot } from 'react-dom/client';
import React, { Component } from 'react';
document.body.innerHTML = '<div id="app"></div>';
import imgUrl from './photos/og-photo.jpg?w=300&h=300'
document.getElementById('meta-image')!.setAttribute("content", imgUrl);
import { Header } from './src/header';
import { ParallaxProvider } from 'react-scroll-parallax';
import { ScrollingBackground } from './src/background';

const root = hydrateRoot(document.getElementById('app'),
	<>
		<ParallaxProvider>
			<ScrollingBackground />
			<Header />
		</ParallaxProvider>
	</>
);
