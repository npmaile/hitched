import TextToSVG from 'text-to-svg';
const textToSVG = TextToSVG.loadSync();

const attributes = { fill: 'red', stroke: 'black' };

const svg = textToSVG.getSVG('L8RSN8R');

console.log(svg);
