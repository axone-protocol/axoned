import { readdir, readFile } from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import { compile } from "@mdx-js/mdx";

const roots = ["docs/command", "docs/predicate", "docs/proto"];

async function* walk(dir) {
    for (const entry of await readdir(dir, { withFileTypes: true })) {
        const entryPath = path.join(dir, entry.name);
        if (entry.isDirectory()) {
            yield* walk(entryPath);
            continue;
        }

        if (entry.isFile() && entry.name.endsWith(".md")) {
            yield entryPath;
        }
    }
}

const failures = [];

for (const root of roots) {
    for await (const file of walk(root)) {
        try {
            await compile(await readFile(file), {
                format: "md",
                filepath: file,
            });
        } catch (error) {
            failures.push({ file, error });
        }
    }
}

if (failures.length > 0) {
    for (const { file, error } of failures) {
        console.error(`MDX compilation failed for ${file}`);
        console.error(error.message);
    }

    process.exit(1);
}
