<template>
	<div>
		<div class="container clearfix" id="main">
			<div class="row">
				<div class="col-3">
					<div class="navigation">
						<div class="accordion accordion-flush" id="accordionFlushExample">
							<div class="accordion-item" id="listdoc">
								<h2 class="accordion-header" id="flush-headingOne">
									<button
										class="accordion-button collapsed"
										type="button"
										data-toggle="collapse"
										data-target="#flush-collapseOne"
										aria-expanded="false"
										aria-controls="flush-collapseOne"
									>
										Documents
									</button>
								</h2>
								<div
									id="flush-collapseOne"
									class="accordion-collapse collapse"
									aria-labelledby="flush-headingOne"
									data-parent="#accordionFlushExample"
								>
									<div class="accordion-body">
										<div class="collection">
											<dl class="content list">
												<dd v-for="(c, i) in listDocument" v-bind:key="i">
													{{ c.group }} -
													<a
														:href="c.href"
														v-on:click="load_document_by_index(i)"
														class="item"
													>
														{{ c.name }}
													</a>
												</dd>
											</dl>
										</div>
									</div>
								</div>
							</div>
							<div class="accordion-item" id="listtoc">
								<h2 class="accordion-header" id="flush-headingTwo">
									<button
										class="accordion-button collapsed"
										type="button"
										data-toggle="collapse"
										data-target="#flush-collapseTwo"
										aria-expanded="false"
										aria-controls="flush-collapseTwo"
									>
										Table Of Contents
									</button>
								</h2>
								<div
									id="flush-collapseTwo"
									class="accordion-collapse collapse"
									aria-labelledby="flush-headingTwo"
									data-parent="#accordionFlushExample"
								>
									<div class="accordion-body">
										<div class="navigation">
											<div class="content list">
												<dl>
													<template v-for="c in listDocumentNavigations">
														<template v-if="c.subgroup.length">
															<dt class="collection title" v-bind:key="c.href">
																<a :href="c.href"> {{ c.name }} </a>
															</dt>

															<dd
																v-for="r in c.subgroup"
																class="content list"
																v-bind:key="r.href"
															>
																<a :href="r.href" class="item">
																	:: {{ r.name }}
																</a>
															</dd>
														</template>
													</template>
												</dl>
											</div>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div class="col-6">
					<div class="markdown markdown-body" v-html="currentDocContent"></div>
				</div>
			</div>
		</div>
		<!-- Back to Top link -->
		<a href="#" class="crunchify-top">↑</a>
	</div>
</template>

<script lang="ts">
import Vue from "vue";
import showdown from "showdown";
import axios from "axios";
import lowerCase from "lodash/lowerCase";

showdown.setOption("tables", true);
showdown.setOption("tasklists", true);
showdown.setOption("ghCodeBlocks", true);
showdown.setOption("metadata", true);
showdown.setOption("simplifiedAutoLink", true);
showdown.setOption("openLinksInNewWindow", true);
showdown.setOption("customizedHeaderId", true);

const converter = new showdown.Converter();

var docfile: any = null;
var constr: any[] = [];

const data = {
	listDocument: constr,
	listDocumentNavigations: constr,
	currentDocContent: "",
	lastrender: Date.now(),
};

const load_startup_document = async () => {
	// fetch document
	const markdownurl =
		"https://raw.githubusercontent.com/avrebarra/postaco/master/WELCOME.md";
	const result = await axios.get(markdownurl);

	// convert document to html
	const docContent = converter.makeHtml(result.data);

	// apply
	data.currentDocContent = docContent;

	// refresh
	await refresh_list_document();
	await refresh_list_document_toc();
	await refresh_content();
};

const load_document_by_index = async (idx: number) => {
	const documentindex = data.listDocument[idx];

	// fetch document
	const markdownurl = documentindex.dochref;
	const result = await axios.get("/" + markdownurl);

	// convert document to html
	const docContent = converter.makeHtml(result.data);

	// apply
	data.currentDocContent = docContent;

	// refresh
	await refresh_list_document();
	await refresh_list_document_toc();
	await refresh_content();
};

const fetch_docfile = async () => {
	const result = await axios.get("/postaco.doc.json");
	docfile = result.data;
};

const refresh_list_document = async () => {
	const idxdocuments: any[] = docfile.documents;
	console.log(idxdocuments[0]);
	data.listDocument as any[];
	data.listDocument = idxdocuments
		.filter((e) => e.kind.includes("document"))
		.map((e) => ({
			group: e.document_dir.substring(1),
			ordering:
				(e.document_dir.substring(1) ? "ZZZZZZZZZZ" : "") +
				e.title.toLowerCase(), // hack to put foldered docs to bottom
			name: e.title,
			dochref: e.document_path_markdown,
		}))
		.sort((a, b) => {
			return a.ordering.toLowerCase().localeCompare(b.ordering.toLowerCase());
		});
};

const refresh_list_document_toc = async () => {
	data.listDocumentNavigations = [];

	var constr: any[] = [];
	let group = {
		name: "",
		href: "",
		subgroup: constr,
	};

	$(".markdown h2, h3").each(function (index, el) {
		const elementType = $(this).prop("nodeName");
		switch (elementType) {
			case "H2":
				if (group.name) data.listDocumentNavigations.push(group); // push first if prev group previously exist

				// create new group
				group = {
					name: $(el).text(),
					href: "#" + $(el).attr("id"),
					subgroup: [],
				};

				break;

			case "H3":
				group.subgroup.push({
					name: $(el).text(),
					href: "#" + $(el).attr("id"),
				});
				break;
		}
	});
};

const refresh_content = async () => {
	$(".markdown img").each(function (index, el) {
		$(el).addClass("img-fluid");
	});

	document.title = docfile.name || docfile.Name;
};

export default Vue.extend({
	methods: {
		load_document_by_index,
	},
	data() {
		return data;
	},
	mounted: async function () {
		await fetch_docfile();
		await refresh_list_document();
		await refresh_list_document_toc();
		await refresh_content();

		await load_startup_document();
	},
});
</script>

<style>
.leftbar {
	background-color: #f9f9fb94;
}

#main.container {
	padding-top: 10vh;
	padding-bottom: 10vh;
}

#listdoc .collection .item {
	cursor: pointer;
}

.navigation {
	max-height: 80vh;
	overflow-y: scroll;
}

.crunchify-top:hover {
	color: #fff !important;
	background-color: #2d2d2d;
	text-decoration: none;
}

.crunchify-top {
	display: none;
	position: fixed;
	bottom: 0rem;
	left: 50vw;
	width: 100vw;
	height: 3.2rem;
	line-height: 3.2rem;
	font-size: 1.4rem;
	color: #fff;
	background-color: rgb(67, 67, 67);
	text-decoration: none;
	text-align: center;
	cursor: pointer;
	transform: translateX(-50%);
}

html {
	scroll-behavior: smooth;
}
</style>
