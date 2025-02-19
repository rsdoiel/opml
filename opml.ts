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

import * as xml from "@libs/xml";

// OPML is the root structure for holding an OPML document
export class OPML {
	version: string = '2.0';
	head: Head = new Head();
	body: Body = new Body();

	// append one or more Body.Outline lists to the current OPML structure
	append(...outlines: OPML[]): boolean {
		let i = this.body.outline.length;
		for (const next of outlines) {
			for (const elem of next.body.outline) {
				this.body.outline.push(elem)
				i += 1;
			}
		}
		if (this.body.outline.length != i) {
			//return fmt.Errorf("Failed to add all outline elements, exlected %d, have %d", i, len(o.Body.Outline))
			return false;
		}
		return true;
	};

	// Walk does a depth first walk of an outline, stops if function
	// return false.
	walk(fn: (outline: Outline) => boolean): boolean {
		if (this.body === undefined || this.body.outline === undefined) {
			return false;
		}
		if (this.body.outline.length > 0) {
			for (const elem of this.body.outline) {
				let ok = walk(elem, fn);
				if (! ok) {
					return false;
				}
			}
			return true;
		}
		return false;
	};

	toObject(): {[key: string]: any} {
		let obj: {[key: string]: any} = {};
		(this.version === '') ? '': obj['@version'] = this.version;
		(this.head === undefined) ? '': obj.head = this.head;
		(this.body === undefined) ? '': obj.body = this.body;
		return obj;
	}

	fromObject(obj: {[key: string]: any}): boolean {
		let ok = true;
		//console.log(`%cDEBUG obj -> ${JSON.stringify(obj)}`, 'color: green')
		if (obj.opml === undefined) {
			(obj.version === undefined ||
				obj.version === '') ? ok = false: this.version = obj.version;
			(obj.head === undefined ||
				Object.keys(obj.head).length === 0) ? ok = false : this.head.fromObject(obj.head);
			(obj.body === undefined ||
				Object.keys(obj.body).length === 0) ? ok = false: this.body = obj.body; 	
		} else {
			(obj.opml.version === undefined ||
				obj.opml === '') ? ok = false: this.version = obj.opml.version;
			(obj.opml.head === undefined ||
				Object.keys(obj.opml.head).length === 0) ? ok = false : this.head.fromObject(obj.opml.head);
			(obj.opml.body === undefined ||
				Object.keys(obj.opml.body).length === 0) ? ok = false: this.body = obj.opml.body;
		}
		console.log(`%cDEBUG this.head -> ${JSON.stringify(this.head.toObject())}`, 'color: magenta')
		return ok;	
	}

	// parse reads a string and returns a OMPL object and error
	parse(src: string): boolean {
		let obj: {[key: string]: any} = {};
		try {
			obj = xml.parse(src);
		} catch (err) {
			return false;
		}
		return this.fromObject(obj);
	}

	// stringify takes the OPML object and returns an XML representation.
	stringify(): string {
		const obj = {
			'@version': '1.0',
			'opml': {
				'@version': '2.0',
				'head': this.head.toObject(),
				'body': this.body.toObject(),
			}
		};
		try {
			return xml.stringify(obj);
		} catch (err) {
			return '';
		}
	}

	async readFile(fname: string): Promise<boolean> {
		const src: string = await Deno.readTextFile(fname);
		try {
			this.parse(src);
		} catch (err) {
			return false;
		}
		return true;
	}

	async writeFile(fName: string): Promise<boolean> {
		const src = this.stringify();
		try {
			await Deno.writeTextFile(fName, src);
		} catch (err) {
			return false;
		}
		return true;
	}
}

// Head holds the metadata for an OPML document
export class Head {
	title: string = '';
	dateCreated: string  = ''; // RFC 882 date and time
	dateModified: string = ''; // RFC 882 date and time
	ownerName: string = '';
	ownerEmail: string = '';
	ownerID: string = '';      // url
	docs: string = '';         // url
	expansionState: string = ''; // array of numbers
	vertScrollState: number = 0;
	windowTop: number = 0;
	windowLeft: number = 0;
	windowBottom: number = 0;
	windowRight: number = 0;

	toObject(): {[key: string]: any} {
		let obj: {[key: string]: any} = {};
		(this.title === '') ? '': obj.title = this.title;
		(this.dateCreated === '') ? '': obj.dateCreated = this.dateCreated;
		(this.dateModified === '') ? '': obj.dateModified = this.dateModified;
		(this.ownerName === '') ? '' : obj.ownerName = this.ownerName;
		(this.ownerEmail === '') ? '': obj.ownerEmail = this.ownerEmail;
		(this.ownerID === '') ? '': obj.ownerID = this.ownerID;
		(this.docs === '') ? '': obj.docs = this.docs;
		(this.expansionState === '') ? '': obj.expansionState = this.expansionState;
		(this.vertScrollState === 0) ? '': obj.vertScrollState = this.vertScrollState;
		(this.windowTop === 0) ? '': obj.windowTop = this.windowTop;
		(this.windowLeft === 0) ? '': obj.windowLeft = this.windowLeft;
		(this.windowBottom === 0) ? '': obj.windowBottom = this.windowBottom;
		(this.windowRight === 0) ? '': obj.windowRight = this.windowRight;
		return obj;
	}

	fromObject(obj: {[key: string]: any}): boolean {
		(obj.title === undefined || obj.title === '') ? '' : this.title = obj.title;
		(obj.dateCreated === undefined || obj.dateCreate === '') ? '': this.dateCreated = obj.dateCreated;
		(obj.dateModified === undefined || obj.dateModified === '') ? '': this.dateModified = obj.dateModified;
		(obj.ownerName === undefined || obj.ownerName === '') ? '' : this.ownerName = obj.ownerName;
		(obj.ownerEmail === undefined || obj.ownerEmail === '') ? '' : this.ownerEmail = obj.ownerEmail;
		(obj.ownerID === undefined || obj.ownerID === '') ? '' : this.ownerID = obj.ownerID;
		(obj.docs === undefined || obj.docs === '') ? '' : this.docs = obj.docs;
		(obj.expansionState === undefined ||
			isNaN(obj.expansionState) || 
			obj.expansionState === 0) ? '' : this.expansionState = obj.expansionState;
		(obj.vertScrollState === undefined ||
			isNaN(obj.vertScrollState) ||
			obj.vertScrollState === 0) ? '' : this.vertScrollState = obj.vertScrollState;
		(obj.windowTop === undefined ||
			isNaN(obj.windowTop) ||
			obj.windowTop === 0) ? '': this.windowTop = obj.windowTop;
		(obj.windowLeft === undefined ||
			isNaN(obj.windowLeft) ||
			obj.windowLeft === 0) ? '' : this.windowLeft = obj.windowLeft;
		(obj.windowBottom === undefined ||
			isNaN(obj.windowBottom) ||
			obj.windowBottom === 0) ? '': this.windowBottom = obj.windowBottom;
		(obj.windowRight === undefined ||
			isNaN(obj.windowRight) ||
			obj.windowRight === 0) ? '': this.windowRight = obj.windowRight;	
		console.log(`%cDEBUG look for title -> ${JSON.stringify(this.toObject())}`, 'color: yellow');
		return true;
	}

}

// Body holds the outline for an OPML document
export class Body {
	outline: Outline[] = [];

	toObject(): {[key: string]: any} {
		let obj:{[key: string]: any} = {};
		(this.outline === undefined ||
			this.outline.length === 0) ? '' : obj.outline = this.outline;
		return obj;
	}

	fromObject(obj: {[key: string]: any}): boolean {
		let ok = true;
		(obj.outline === undefined ||
			obj.outline.length === 0) ? ok = false : this.outline = obj.outline;
		return ok;	
	}
}


// Outline is the primary element of an OPML document, may hold sub-Outlines
export class Outline {
	text: string = '';
	type: string = '';
	title: string = '';
	isComment: boolean = false;
	isBreakpoint: boolean = false;
	created: string = ''; // RFC 882 date and time
	category: string = '';
	xmlURL: string = '';  // url
	htmlURL: string = ''; // url
	language: string = '';
	description: string = '';
	version: string = '';
	url: string = '';     // url
	outline: Outline[] = [];
	otherAttr: {[key: string]: any} = {};

	// HasChildren return true if the outline element has a populated child outline
	hasChildren(): boolean {
		if (this.outline !== undefined && this.outline.length > 0) {
			return true;
		}
		return false;
	};
}

// New creates an empty OPML structure
export function NewOPML(): OPML {
	let o = new OPML();
	// OPML spec support
	o.version = `2.0`;
	o.head = new Head();
	o.body = new Body();
	return o
}

// ReadFile reads a OPML file and returns a new OPML structure and error
export async function readFile(fname: string): Promise<OPML | undefined> {
	const o = NewOPML();
	try {
		await o.readFile(fname);
	} catch (err) {
		return undefined;
	}
	return o;
}

// This is a general walk function for outline.
export function walk(ol: Outline, fn: (outLine: Outline) => boolean): boolean {
	if (ol === undefined)  {
		return false;
	}
	const ok = fn(ol);
	if (ol.outline !== undefined) {
		for (const elem of ol.outline) {
			return walk(elem, fn);
		}
	}
	return ok;
};
