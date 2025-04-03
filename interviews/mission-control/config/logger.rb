require "semantic_logger"
level = ENV.fetch("LOG_LEVEL", "error").to_sym
SemanticLogger.default_level = level
SemanticLogger.add_appender(io: $stdout, level: level)
