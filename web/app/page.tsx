'use client';
import { useState } from 'react';
import Link from 'next/link';
import {
  ClipboardDocumentIcon,
  ChartBarIcon,
  SparklesIcon,
  PencilSquareIcon,
  CommandLineIcon,
  ArrowDownTrayIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
} from '@heroicons/react/24/outline';

interface ReviewIssue {
  type: string;
  file?: string;
  line?: number;
  message: string;
  severity: string;
}

interface ReviewRecord {
  id: string;
  pull_request: string;
  title: string;
  reviewer: string;
  summary: string;
  issues: ReviewIssue[];
  score: number;
  created_at: string;
}

export default function Home() {
  const [prUrl, setPrUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [review, setReview] = useState<any>(null);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);

  const handleReview = async () => {
    if (!prUrl) return;
    setLoading(true);
    setError('');
    setReview(null);
    
    try {
      const res = await fetch('/api/review', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pr_url: prUrl }),
      });
      
      if (!res.ok) {
        throw new Error('Failed to generate review');
      }
      
      const data = await res.json();
      setReview(data.result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async () => {
    if (!review) return;
    const text = `${review.summary}\n\nIssues:\n${review.issues?.map((i: ReviewIssue) => `- ${i.message}`).join('\n')}`;
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'high': return 'text-red-400 bg-red-900/20 border-red-500/30';
      case 'medium': return 'text-yellow-400 bg-yellow-900/20 border-yellow-500/30';
      default: return 'text-blue-400 bg-blue-900/20 border-blue-500/30';
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <div className="max-w-4xl mx-auto px-4 py-12">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-5xl font-bold text-white mb-4 flex items-center justify-center gap-3">
            <SparklesIcon className="w-12 h-12 text-purple-400" />
            AI Code Review
          </h1>
          <p className="text-xl text-purple-200">Automated PR analysis and feedback</p>
        </div>

        {/* Input */}
        <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-8 mb-8 border border-purple-500/20">
          <div className="flex gap-4">
            <input
              type="text"
              placeholder="Paste GitHub PR URL here..."
              value={prUrl}
              onChange={(e) => setPrUrl(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleReview()}
              className="flex-1 px-6 py-4 bg-slate-900/50 border border-purple-500/30 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:border-purple-500 text-lg"
            />
            <button
              onClick={handleReview}
              disabled={loading || !prUrl}
              className="px-8 py-4 bg-purple-600 hover:bg-purple-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-xl transition-all font-semibold text-lg flex items-center gap-2"
            >
              <ClipboardDocumentIcon className="w-6 h-6" />
              {loading ? 'Reviewing...' : 'Review'}
            </button>
          </div>
          {error && <p className="mt-4 text-red-400">{error}</p>}
        </div>

        {/* Results */}
        {review && (
          <div className="space-y-6">
            {/* Score Card */}
            <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
              <div className="flex justify-between items-center">
                <div>
                  <h2 className="text-2xl font-bold text-white mb-2">{review.title || 'Code Review'}</h2>
                  <div className="flex gap-4 text-gray-400">
                    <span>{review.source || 'AI Review'}</span>
                    <span>•</span>
                    <span>{new Date().toLocaleDateString()}</span>
                  </div>
                </div>
                <div className="text-right">
                  <div className={`text-4xl font-bold ${review.score >= 80 ? 'text-green-400' : review.score >= 60 ? 'text-yellow-400' : 'text-red-400'}`}>
                    {review.score.toFixed(0)}
                  </div>
                  <div className="text-gray-400 text-sm">Quality Score</div>
                </div>
              </div>
            </div>

            {/* Summary */}
            <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-xl font-semibold text-white flex items-center gap-2">
                  <PencilSquareIcon className="w-6 h-6 text-purple-400" />
                  Summary
                </h3>
                <button
                  onClick={copyToClipboard}
                  className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg text-sm transition-all flex items-center gap-2"
                >
                  {copied ? <CheckCircleIcon className="w-4 h-4" /> : <ClipboardDocumentIcon className="w-4 h-4" />}
                  {copied ? 'Copied!' : 'Copy'}
                </button>
              </div>
              <p className="text-gray-300 leading-relaxed">{review.summary}</p>
            </div>

            {/* Issues */}
            {review.issues && review.issues.length > 0 && (
              <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
                <h3 className="text-xl font-semibold text-white mb-4 flex items-center gap-2">
                  <ExclamationTriangleIcon className="w-6 h-6 text-yellow-400" />
                  Issues ({review.issues.length})
                </h3>
                <div className="space-y-3">
                  {review.issues.map((issue: ReviewIssue, idx: number) => (
                    <div key={idx} className={`p-4 rounded-lg border ${getSeverityColor(issue.severity)}`}>
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <span className="text-xs font-medium uppercase tracking-wide">{issue.type}</span>
                          {issue.file && (
                            <span className="ml-3 text-xs text-gray-400">
                              {issue.file}{issue.line ? `:${issue.line}` : ''}
                            </span>
                          )}
                          <p className="mt-1">{issue.message}</p>
                        </div>
                        <span className={`px-2 py-1 rounded text-xs font-medium ml-4 ${
                          issue.severity === 'high' ? 'bg-red-500/20 text-red-300' :
                          issue.severity === 'medium' ? 'bg-yellow-500/20 text-yellow-300' :
                          'bg-blue-500/20 text-blue-300'
                        }`}>
                          {issue.severity}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Suggestions */}
            {review.suggestions && review.suggestions.length > 0 && (
              <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-purple-500/20">
                <h3 className="text-xl font-semibold text-white mb-4 flex items-center gap-2">
                  <ArrowDownTrayIcon className="w-6 h-6 text-green-400" />
                  Suggestions
                </h3>
                <ul className="space-y-2">
                  {review.suggestions.map((s: string, idx: number) => (
                    <li key={idx} className="flex items-start gap-3 text-gray-300">
                      <span className="w-2 h-2 bg-green-400 rounded-full mt-2 flex-shrink-0"></span>
                      {s}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        )}

        {/* Features */}
        {!review && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <ClipboardDocumentIcon className="w-12 h-12 text-purple-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">PR Analysis</h3>
              <p className="text-gray-400">Automatically extract and analyze code changes from GitHub PRs</p>
            </div>
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <SparklesIcon className="w-12 h-12 text-yellow-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">AI Review</h3>
              <p className="text-gray-400">Get intelligent feedback on code quality, security, and best practices</p>
            </div>
            <div className="bg-slate-800/30 rounded-xl p-6 text-center">
              <ArrowDownTrayIcon className="w-12 h-12 text-green-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-white mb-2">Export Reports</h3>
              <p className="text-gray-400">Download review reports in Markdown, JSON, or plain text formats</p>
            </div>
          </div>
        )}

        {/* Quick Links */}
        <div className="flex justify-center gap-4 mt-8">
          <Link href="/history" className="px-6 py-3 bg-slate-800/50 hover:bg-slate-700/50 text-white rounded-xl transition-all flex items-center gap-2">
            <ClipboardDocumentIcon className="w-5 h-5" />
            History
          </Link>
          <Link href="/dashboard" className="px-6 py-3 bg-slate-800/50 hover:bg-slate-700/50 text-white rounded-xl transition-all flex items-center gap-2">
            <ChartBarIcon className="w-5 h-5" />
            Analytics
          </Link>
        </div>
      </div>
    </div>
  );
}
