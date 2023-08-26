const { readFile, writeFile, readDir, toYaml, fromYaml, scan } = require('nuv');

// test reading a file
let data = readFile('testdata/sample.txt');
console.log(data);
console.log('***');

// test writing a file
data = 'test write';
writeFile('testdata/sample2.txt', data);

data = readFile('testdata/sample2.txt');
console.log(data);
console.log('***');

// test converting to yaml
let yaml = toYaml({ a: 1, b: 2 });
console.log(yaml);
console.log('***');

// test converting from yaml
let obj = fromYaml(yaml);
console.log(JSON.stringify(obj));
console.log('***');

// test scanning a directory
let scanResults = scan('testdata', (folder) => folder + ' ');
console.log(scanResults);
console.log('***');

// test reading a directory
let files = readDir('testdata');
console.log(files);