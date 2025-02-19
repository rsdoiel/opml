//
// Package opml provides basic utility functions for working with OPML files.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
//
// Copyright (c) 2025, R. S. Doiel
// All rights not granted herein are expressly reserved by R. S. Doiel.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
import { assertEquals, assertNotEquals } from '@std/assert';
import * as path from "@std/path";
import * as opml from './opml.ts';

Deno.test("Test NewOPML and stringify", 
function TestNewAndString() {
	let o = opml.NewOPML();
	assertNotEquals(o, undefined, "Can't create an opml structure");

	let head = o.head;
	assertNotEquals(head, undefined, "Can't find Head in opml structure");

	let body = o.body;
	assertNotEquals(body, undefined, "Can't find Body in opml structure");
	o.body.outline = [];
	
	let expected = `<?xml version="1.0"?>
<opml version="2.0">
  <head/>
  <body/>
</opml>`;
	let src = o.stringify();
	assertEquals(expected, src, `expected ${expected}, got ${src}`);
});

Deno.test("Test Read", 
async function TestRead() {
	let o = opml.NewOPML();
	let fname = "testdata/example1.opml";
	assertEquals(await o.readFile(fname), true,
		`o.readFile() should return an OPML structure, got ${o.stringify()}`);
	assertEquals(o.version, "2.0", `Expected version 2.0, got ${o.version}`);
	assertNotEquals(o.head, undefined, "expected an head element.");
	assertNotEquals(o.head.title, undefined, "expected an head.title element");
	assertNotEquals(o.head.title, '', `expected a populated head.title element, read may have failed to map data from ${fname}, ${JSON.stringify(o.head.toObject())}`);
	assertEquals(o.head.title,
		"johndoe@example.com subscriptions in Go Read",
		`Expected "johndoe@example.com subscriptions in Go Read", found -> '${JSON.stringify(o.head)}' in ${fname}`);
	let i = 64;
	assertEquals(o.body.outline.length, i, 
		`expected ${i} outline elements, found ${o.body.outline.length}`);
	let s = o.stringify();
	assertNotEquals(s.indexOf(`<outline text=""></outline>`), -1,
		`an empty outline is included in string: ${s}`);

	o = opml.NewOPML();
	assertEquals(await o.readFile("testdata/example2.opml"), true, 
		`ReadFile should return an OPML structure`)
	assertEquals(o.version, "2.0", 
		`Expected version 2.0, got ${o.version}`);
});

Deno.test("Test Write", async function TestWrite() {
	const fname = path.join("testdata", "test1.opml")
	let o = opml.NewOPML();
	let txt = o.stringify();
	await o.writeFile(fname);

	let src = await Deno.readTextFile(fname);
	assertNotEquals(src, undefined, `should have been able to reaad ${fname}`);
	assertEquals(txt, src, `expected\n\t${txt}\n, got\n\t${src}\n`);
	// cleanup the temp file
	//await Deno.remove(fname);
});

Deno.test("Test Append", async function TestAppend() {
	let o1 = opml.NewOPML();
	let fname = "testdata/simple1.opml";
	assertEquals(await o1.readFile("testdata/simple1.opml"), true, `failed to read ${fname}`);
	let o2 = opml.NewOPML();
	fname ="testdata/simple2.opml";
	assertEquals(await o2.readFile(fname), true, `failed to read ${fname}`);

	let o3 = opml.NewOPML();
	assertEquals(o3.append(o1, o2), true, `expected o3 to be able to append o1 and o2`);

	let tot = o1.body.outline.length + o2.body.outline.length;
	assertEquals(o3.body.outline.length, tot,
		`expected ${tot} outline elements, got ${o3.body.outline.length}`);
});

