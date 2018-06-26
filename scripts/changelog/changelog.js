const Q = require('q');
const path = require('path');
const readFile = Q.denodeify(require(`fs`).readFile);

const parserOpts = {
	headerPattern: /^:(\w*): (.*)$/,
	headerCorrespondence: [`type`, `subject`],
	noteKeywords: [`BREAKING CHANGE`, `BREAKING CHANGES`],
};

const templateDir = path.dirname(require.resolve('conventional-changelog-emoji'));
const writerOpts = Q.all(
    ['template', 'header', 'commit', 'footer']
        .map(name => readFile(path.resolve(templateDir, `./templates/${name}.hbs`), `utf-8`))
).spread((template, header, commit, footer) => ({
    transform: (commit, context) => {
        let discard = true;
        const issues = [];

        commit.notes.forEach(note => {
            note.title = `BREAKING CHANGES`;
            discard = false;
        });

        switch(commit.type) {
            case 'sparkles':
                commit.group = "<!-- 1 -->:sparkles: Features";
                break;
            case 'bug':
                commit.group = "<!-- 2 -->:bug: Bugfixes";
                break;
            case 'zap':
            case 'lock':
            case 'lipstick':
            case 'recycle':
                commit.group = "<!-- 3 -->:tada: Improvements"
                break;
            case 'whale':
            case 'construction_worker':
            case 'rocket':
            case 'wrench':
                commit.group = "<!-- 4 -->:wrench: Tooling"
                break;
            default:
                commit.group = "Other commits";
                break;
        }

        if (typeof commit.hash === `string`) {
            commit.hash = commit.hash.substring(0, 7);
        }

        if (typeof commit.subject === `string`) {
            let url = context.repository
                ? `${context.host}/${context.owner}/${context.repository}`
                : context.repoUrl;
            if (url) {
                url = `${url}/issues/`;
                // Issue URLs.
                commit.subject = commit.subject.replace(
                    /#([0-9]+)/g,
                    (_, issue) => {
                        issues.push(issue);
                        return `[#${issue}](${url}${issue})`;
                    },
                );
            }
            if (context.host) {
                // User URLs.
                commit.subject = commit.subject.replace(
                    /\B@([a-z0-9](?:-?[a-z0-9]){0,38})/g,
                    `[@$1](${context.host}/$1)`,
                );
            }
        }

        // remove references that already appear in the subject
        commit.references = commit.references.filter(reference => {
            if (issues.indexOf(reference.issue) === -1) {
                return true;
            }

            return false;
        });

        return commit;
    },
    groupBy: `group`,
    commitGroupsSort: `title`,
    commitsSort: `committerDate`,
    noteGroupsSort: `title`,
    mainTemplate: template,
    headerPartial: header,
    commitPartial: commit,
    footerPartial: footer,
}));

module.exports = Q.all([
	parserOpts,
	writerOpts,
]).spread(
	(parserOpts, writerOpts) => {
		return {
			parserOpts,
			writerOpts,
		};
	},
);
