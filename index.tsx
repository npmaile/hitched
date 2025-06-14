import { hydrateRoot } from 'react-dom/client';
import { Component } from 'react';
document.body.innerHTML = '<div id="app"></div>';
import imgUrl from './photos/og-photo.jpg?w=300&h=300'
document.getElementById('meta-image')!.setAttribute("content",imgUrl);
import { Header } from './src/header';
import {ParallaxProvider} from 'react-scroll-parallax';
import {ScrollingBackground} from './src/background';


class Abc extends Component {
	render(){
		return <div class="text-5xl font-bold">information will come soon</div>
	}
}

const root = hydrateRoot(document.getElementById('app'),
			 	<>
<ParallaxProvider>
<ScrollingBackground>
</ScrollingBackground>
				<Header></Header>
			 	<h1>
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					Nate and Susy wedding page
					<Abc/>
					<img src={imgUrl}></img>
				</h1>
				</ParallaxProvider>
				</>
			);
