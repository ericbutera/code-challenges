# frozen_string_literal: true

# Game statistics
class Stats < Struct.new(
  :abort_count,
  :explode_count,
  :retry_count,
  :distance_traveled_km,
  :fuel_burned_liters,
  :flight_time_seconds,
  :mission_count
)
  def initialize(attrs = {})
    super(*members.map { |field| attrs[field] || 0 })
  end

  def merge(other)
    each_pair do |name, _|
      self[name] += other[name]
    end
  end

  def ==(other)
    return false unless other.is_a?(Stats)

    each_pair.all? { |name, value| value == other[name] }
  end
end
